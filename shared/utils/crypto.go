package utils

import (
	"crypto/rand"
	"encoding/binary"
)

func CryptoRandInt64() (int64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	n := int64(binary.BigEndian.Uint64(b[:]))
	return n, nil
}
