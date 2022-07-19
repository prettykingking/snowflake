// Package snowflake implements Snowflake, a distributed unique ID generator inspired by Twitter's Snowflake.
//
// A Snowflake ID is composed of
//     39 bits for time in units of 10 msec
//      8 bits for a sequence number
//     16 bits for a machine id
package snowflake

import (
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

// These constants are the bit lengths of Snowflake ID parts.
const (
	BitLenTime      = 39                               // bit length of time
	BitLenSequence  = 8                                // bit length of sequence number
	BitLenMachineId = 63 - BitLenTime - BitLenSequence // bit length of machine id
)

// Settings configures Snowflake:
//
// StartTime is the time since which the Snowflake time is defined as the elapsed time.
// If StartTime is 0, the start time of the Snowflake is set to "2014-09-01 00:00:00 +0000 UTC".
// If StartTime is ahead of the current time, Snowflake is not created.
//
// MachineId returns the unique ID of the Snowflake instance.
// If MachineId returns an error, Snowflake is not created.
// If MachineId is nil, default MachineId is used.
// Default MachineId returns the lower 16 bits of the private IP address.
//
// CheckMachineId validates the uniqueness of the machine ID.
// If CheckMachineId returns false, Snowflake is not created.
// If CheckMachineId is nil, no validation is done.
type Settings struct {
	StartTime      time.Time
	MachineId      func() (uint16, error)
	CheckMachineId func(uint16) bool
}

// Snowflake is a distributed unique ID generator.
type Snowflake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineId   uint16
}

// NewSnowflake returns a new Snowflake configured with the given Settings.
// NewSnowflake returns nil in the following cases:
// - Settings.StartTime is ahead of the current time.
// - Settings.MachineId returns an error.
// - Settings.CheckMachineId returns false.
func NewSnowflake(st Settings) *Snowflake {
	sf := new(Snowflake)
	sf.mutex = new(sync.Mutex)
	sf.sequence = uint16(1<<BitLenSequence - 1)

	if st.StartTime.After(time.Now()) {
		return nil
	}
	if st.StartTime.IsZero() {
		return nil
	} else {
		sf.startTime = toSnowflakeTime(st.StartTime)
	}

	var err error
	if st.MachineId == nil {
		return nil
	} else {
		sf.machineId, err = st.MachineId()
	}
	if err != nil || (st.CheckMachineId != nil && !st.CheckMachineId(sf.machineId)) {
		return nil
	}

	return sf
}

// NextID generates a next unique ID.
// After the Snowflake time overflows, NextID returns an error.
func (sf *Snowflake) NextID() (uint64, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime(overtime))
		}
	}

	return sf.toID()
}

const snowflakeTimeUnit = 1e7 // nsec, i.e. 10 msec

func toSnowflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / snowflakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toSnowflakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%snowflakeTimeUnit)*time.Nanosecond
}

func (sf *Snowflake) toID() (uint64, error) {
	if sf.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("over the time limit")
	}

	return uint64(sf.elapsedTime)<<(BitLenSequence+BitLenMachineId) |
		uint64(sf.sequence)<<BitLenMachineId |
		uint64(sf.machineId), nil
}

func (sf *Snowflake) Marshal() string {
	var id [5]byte
	// time in 10 ms
	binary.BigEndian.PutUint64(id[:], uint64(sf.elapsedTime))
	// sequence

	return ""
}

// Decompose returns a set of Snowflake ID parts.
func Decompose(id uint64) map[string]uint64 {
	const maskSequence = uint64((1<<BitLenSequence - 1) << BitLenMachineId)
	const maskMachineID = uint64(1<<BitLenMachineId - 1)

	msb := id >> 63
	t := id >> (BitLenSequence + BitLenMachineId)
	sequence := id & maskSequence >> BitLenMachineId
	machineId := id & maskMachineID
	return map[string]uint64{
		"id":        id,
		"msb":       msb,
		"time":      t,
		"sequence":  sequence,
		"machineId": machineId,
	}
}
