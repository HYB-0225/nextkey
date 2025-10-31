package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	CharsetLetters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetAlphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

const (
	CharsetTypeLetters      = "letters"
	CharsetTypeAlphanumeric = "alphanumeric"
)

func getCharset(charsetType string) string {
	switch charsetType {
	case CharsetTypeLetters:
		return CharsetLetters
	case CharsetTypeAlphanumeric:
		return CharsetAlphanumeric
	default:
		return CharsetAlphanumeric
	}
}

func RandomString(length int, charsetType string) string {
	charset := getCharset(charsetType)
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func GenerateCardKey(prefix, suffix string, length int, charsetType string) string {
	if length < 6 {
		length = 6
	}
	if length > 32 {
		length = 32
	}
	return prefix + RandomString(length, charsetType) + suffix
}
