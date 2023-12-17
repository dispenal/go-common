package kafka

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func hashMessage(btext []byte) string {
	hash := md5.Sum(btext)
	return hex.EncodeToString(hash[:])
}
