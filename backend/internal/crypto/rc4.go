package crypto

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func init() {
	Register(EncryptorFactory{
		Meta: EncryptorMeta{
			Scheme:        "rc4",
			Name:          "RC4",
			Description:   "传统流加密算法，仅用于兼容性需求",
			SecurityLevel: "insecure",
			Performance:   "fast",
			IsDeprecated:  true,
		},
		NewEncryptor: func(key string) (Encryptor, error) {
			return NewRC4Encryptor(key)
		},
		GenerateKey: generateRC4Key,
	})
}

type RC4Encryptor struct {
	key []byte
}

func NewRC4Encryptor(key string) (*RC4Encryptor, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		// 尝试直接使用字符串
		keyBytes = []byte(key)
	}

	if len(keyBytes) == 0 {
		return nil, errors.New("RC4密钥不能为空")
	}

	return &RC4Encryptor{key: keyBytes}, nil
}

func (e *RC4Encryptor) Scheme() string {
	return "rc4"
}

func generateRC4Key() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (e *RC4Encryptor) Encrypt(plaintext string) (string, error) {
	cipher, err := rc4.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plaintext))
	cipher.XORKeyStream(ciphertext, []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *RC4Encryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	cipher, err := rc4.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(data))
	cipher.XORKeyStream(plaintext, data)

	return string(plaintext), nil
}
