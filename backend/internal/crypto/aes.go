package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var globalKey []byte

func SetKey(key string) error {
	keyBytes, err := decodeKey(key)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return errors.New("AES key must be 32 bytes")
	}
	globalKey = keyBytes
	return nil
}

func decodeKey(key string) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err == nil && len(keyBytes) == 32 {
		return keyBytes, nil
	}
	if len(key) == 64 {
		return []byte(key)[:32], nil
	}
	return []byte(key), nil
}

func Encrypt(plaintext string) (string, error) {
	if globalKey == nil {
		return "", errors.New("encryption key not set")
	}

	block, err := aes.NewCipher(globalKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) (string, error) {
	if globalKey == nil {
		return "", errors.New("encryption key not set")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(globalKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
