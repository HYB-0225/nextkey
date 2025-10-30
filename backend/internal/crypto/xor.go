package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func init() {
	Register(EncryptorFactory{
		Meta: EncryptorMeta{
			Scheme:        "xor",
			Name:          "XOR",
			Description:   "简单异或加密，仅用于兼容性需求",
			SecurityLevel: "insecure",
			Performance:   "fast",
			IsDeprecated:  true,
		},
		NewEncryptor: func(key string) (Encryptor, error) {
			return NewXOREncryptor(key)
		},
		GenerateKey: generateXORKey,
	})
}

type XOREncryptor struct {
	key []byte
}

func NewXOREncryptor(key string) (*XOREncryptor, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		// 尝试直接使用字符串
		keyBytes = []byte(key)
	}

	if len(keyBytes) == 0 {
		return nil, errors.New("XOR密钥不能为空")
	}

	return &XOREncryptor{key: keyBytes}, nil
}

func (e *XOREncryptor) Scheme() string {
	return "xor"
}

func generateXORKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (e *XOREncryptor) Encrypt(plaintext string) (string, error) {
	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, len(plaintextBytes))

	keyLen := len(e.key)
	for i := 0; i < len(plaintextBytes); i++ {
		ciphertext[i] = plaintextBytes[i] ^ e.key[i%keyLen]
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *XOREncryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(data))

	keyLen := len(e.key)
	for i := 0; i < len(data); i++ {
		plaintext[i] = data[i] ^ e.key[i%keyLen]
	}

	return string(plaintext), nil
}
