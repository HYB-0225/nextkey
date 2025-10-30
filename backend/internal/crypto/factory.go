package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func NewEncryptor(scheme, key string) (Encryptor, error) {
	switch scheme {
	case "aes-256-gcm":
		return NewAESEncryptor(key)
	default:
		return nil, errors.New("unsupported encryption scheme")
	}
}

func GenerateEncryptionKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
