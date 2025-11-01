package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

func init() {
	Register(EncryptorFactory{
		Meta: EncryptorMeta{
			Scheme:        "chacha20-poly1305",
			Name:          "ChaCha20-Poly1305",
			Description:   "现代高性能AEAD加密算法，移动端友好",
			SecurityLevel: "secure",
			Performance:   "fast",
			IsDeprecated:  false,
		},
		NewEncryptor: func(key string) (Encryptor, error) {
			return NewChaCha20Encryptor(key)
		},
		GenerateKey: generateChaCha20Key,
	})
}

type ChaCha20Encryptor struct {
	key []byte
}

func NewChaCha20Encryptor(key string) (*ChaCha20Encryptor, error) {
	keyBytes, err := decodeKey(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != chacha20poly1305.KeySize {
		return nil, errors.New("ChaCha20密钥必须为32字节")
	}
	return &ChaCha20Encryptor{key: keyBytes}, nil
}

func (e *ChaCha20Encryptor) Scheme() string {
	return "chacha20-poly1305"
}

func generateChaCha20Key() string {
	bytes := make([]byte, chacha20poly1305.KeySize)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (e *ChaCha20Encryptor) Encrypt(plaintext string) (string, error) {
	aead, err := chacha20poly1305.New(e.key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *ChaCha20Encryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	aead, err := chacha20poly1305.New(e.key)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
