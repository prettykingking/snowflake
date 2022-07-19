package snowflake

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"testing"
	"time"
)

var sf *Snowflake

var startTime int64
var machineId uint64

func init() {
	ip, _ := lower16BitPrivateIP()
	machineId = uint64(ip)

	var st Settings
	st.StartTime = time.Now()
	st.MachineId = func() (uint16, error) {
		return ip, nil
	}

	startTime = toSnowflakeTime(st.StartTime)

	sf = NewSnowflake(st)
	if sf == nil {
		panic("snowflake not created")
	}
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

func nextID(t *testing.T) uint64 {
	id, err := sf.NextID()
	if err != nil {
		t.Fatal("id not generated")
	}
	return id
}

func TestSnowflakeOnce(t *testing.T) {
	sleepTime := uint64(50)
	time.Sleep(time.Duration(sleepTime) * 10 * time.Millisecond)

	id := nextID(t)
	parts := Decompose(id)

	actualMSB := parts["msb"]
	if actualMSB != 0 {
		t.Errorf("unexpected msb: %d", actualMSB)
	}

	actualTime := parts["time"]
	if actualTime < sleepTime || actualTime > sleepTime+1 {
		t.Errorf("unexpected time: %d", actualTime)
	}

	actualSequence := parts["sequence"]
	if actualSequence != 0 {
		t.Errorf("unexpected sequence: %d", actualSequence)
	}

	actualMachineID := parts["machineId"]
	if actualMachineID != machineId {
		t.Errorf("unexpected machine id: %d", actualMachineID)
	}

	fmt.Println("snowflake id:", id)
	fmt.Println("decompose:", parts)
}

func currentTime() int64 {
	return toSnowflakeTime(time.Now())
}

func TestEncoding(t *testing.T) {
	id := nextID(t)

	idHex := ToHex(id)
	if len(idHex) != 13 {
		t.Errorf("corrupted id encoding")
	}

	idInt, err := ToInt(idHex)
	if err != nil {
		t.Error(err)
	}

	if idInt != id {
		t.Errorf("corrupted id decoding")
	}
}

func TestSnowflakeFor10Sec(t *testing.T) {
	var numID uint32
	var lastID uint64
	var maxSequence uint64

	initial := currentTime()
	current := initial
	for current-initial < 1000 {
		id := nextID(t)
		parts := Decompose(id)
		numID++

		if id <= lastID {
			t.Fatal("duplicated id")
		}
		lastID = id

		current = currentTime()

		actualMSB := parts["msb"]
		if actualMSB != 0 {
			t.Errorf("unexpected msb: %d", actualMSB)
		}

		actualTime := int64(parts["time"])
		overtime := startTime + actualTime - current
		if overtime > 0 {
			t.Errorf("unexpected overtime: %d", overtime)
		}

		actualSequence := parts["sequence"]
		if maxSequence < actualSequence {
			maxSequence = actualSequence
		}

		actualMachineID := parts["machineId"]
		if actualMachineID != machineId {
			t.Errorf("unexpected machine id: %d", actualMachineID)
		}
	}

	if maxSequence != 1<<BitLenSequence-1 {
		t.Errorf("unexpected max sequence: %d", maxSequence)
	}
	fmt.Println("max sequence:", maxSequence)
	fmt.Println("number of id:", numID)
}

func TestSnowflakeInParallel(t *testing.T) {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Println("number of cpu:", numCPU)

	consumer := make(chan uint64)

	const numID = 10000
	generate := func() {
		for i := 0; i < numID; i++ {
			consumer <- nextID(t)
		}
	}

	const numGenerator = 10
	for i := 0; i < numGenerator; i++ {
		go generate()
	}

	set := make(map[uint64]struct{})
	for i := 0; i < numID*numGenerator; i++ {
		id := <-consumer
		if _, ok := set[id]; ok {
			t.Fatal("duplicated id")
		}
		set[id] = struct{}{}
	}

	fmt.Println("number of id:", len(set))
}

func TestNilSnowflake(t *testing.T) {
	var startInFuture Settings
	startInFuture.StartTime = time.Now().Add(time.Duration(1) * time.Minute)
	if NewSnowflake(startInFuture) != nil {
		t.Errorf("snowflake starting in the future")
	}

	var noMachineID Settings
	noMachineID.MachineId = func() (uint16, error) {
		return 0, fmt.Errorf("no machine id")
	}
	if NewSnowflake(noMachineID) != nil {
		t.Errorf("snowflake with no machine id")
	}

	var invalidMachineID Settings
	invalidMachineID.CheckMachineId = func(uint16) bool {
		return false
	}
	if NewSnowflake(invalidMachineID) != nil {
		t.Errorf("snowflake with invalid machine id")
	}
}

func pseudoSleep(period time.Duration) {
	sf.startTime -= int64(period) / snowflakeTimeUnit
}

func TestNextIDError(t *testing.T) {
	year := time.Duration(365*24) * time.Hour
	pseudoSleep(time.Duration(174) * year)
	nextID(t)

	pseudoSleep(time.Duration(1) * year)
	_, err := sf.NextID()
	if err == nil {
		t.Errorf("time is not over")
	}
}
