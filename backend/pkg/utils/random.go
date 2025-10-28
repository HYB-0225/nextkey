package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func GenerateCardKey(prefix, suffix string, length int) string {
	if length < len(prefix)+len(suffix) {
		length = len(prefix) + len(suffix) + 8
	}
	randomLen := length - len(prefix) - len(suffix)
	return prefix + RandomString(randomLen) + suffix
}
