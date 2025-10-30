package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

func init() {
	Register(EncryptorFactory{
		Meta: EncryptorMeta{
			Scheme:        "aes-256-gcm",
			Name:          "AES-256-GCM",
			Description:   "高强度对称加密算法，推荐使用",
			SecurityLevel: "secure",
			IsDeprecated:  false,
		},
		NewEncryptor: func(key string) (Encryptor, error) {
			return NewAESEncryptor(key)
		},
		GenerateKey: generateAESKey,
	})
}

type AESEncryptor struct {
	key []byte
}

func NewAESEncryptor(key string) (*AESEncryptor, error) {
	keyBytes, err := decodeKey(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, errors.New("AES key must be 32 bytes")
	}
	return &AESEncryptor{key: keyBytes}, nil
}

func (e *AESEncryptor) Scheme() string {
	return "aes-256-gcm"
}

func generateAESKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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

func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
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

func (e *AESEncryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
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
