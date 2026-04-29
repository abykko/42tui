package server

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecret() string {
	b := make([]byte, 32) // 256 bits

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}