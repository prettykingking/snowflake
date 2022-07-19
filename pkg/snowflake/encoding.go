package snowflake

import (
	"encoding/base32"
	"encoding/binary"
)

const encodeHex = "0123456789abcdefghijklmnopqrstuv"

func ToHex(id uint64) string {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], id)

	enc := base32.NewEncoding(encodeHex).WithPadding(base32.NoPadding)

	return enc.EncodeToString(b[:])
}

func ToInt(id string) (uint64, error) {
	enc := base32.NewEncoding(encodeHex).WithPadding(base32.NoPadding)
	var dst [8]byte
	_, err := enc.Decode(dst[:], []byte(id))
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(dst[:]), nil
}
