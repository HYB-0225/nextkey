package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func init() {
	Register(EncryptorFactory{
		Meta: EncryptorMeta{
			Scheme:        "custom-base64",
			Name:          "自定义Base64",
			Description:   "随机字符表Base64编码，简单混淆",
			SecurityLevel: "insecure",
			Performance:   "fast",
			IsDeprecated:  false,
		},
		NewEncryptor: func(key string) (Encryptor, error) {
			return NewCustomBase64Encryptor(key)
		},
		GenerateKey: generateCustomBase64Key,
	})
}

type CustomBase64Encryptor struct {
	encoding *base64.Encoding
}

func NewCustomBase64Encryptor(key string) (*CustomBase64Encryptor, error) {
	if len(key) != 64 {
		return nil, errors.New("自定义Base64密钥必须为64个字符")
	}

	// 验证字符表是否有效（64个不同字符）
	charSet := make(map[rune]bool)
	for _, c := range key {
		if charSet[c] {
			return nil, errors.New("自定义Base64密钥包含重复字符")
		}
		charSet[c] = true
	}

	if len(charSet) != 64 {
		return nil, errors.New("自定义Base64密钥必须包含恰好64个不同字符")
	}

	encoding := base64.NewEncoding(key)

	return &CustomBase64Encryptor{encoding: encoding}, nil
}

func (e *CustomBase64Encryptor) Scheme() string {
	return "custom-base64"
}

func generateCustomBase64Key() string {
	// 标准Base64字符集
	standardChars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	// Fisher-Yates洗牌算法
	chars := make([]rune, len(standardChars))
	copy(chars, standardChars)

	for i := len(chars) - 1; i > 0; i-- {
		// 生成随机索引
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		j := int(randomBytes[0]) % (i + 1)

		// 交换
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

func (e *CustomBase64Encryptor) Encrypt(plaintext string) (string, error) {
	// 使用自定义字符表编码
	customEncoded := e.encoding.EncodeToString([]byte(plaintext))

	// 外层再用标准Base64包裹，避免SDK解析问题
	return base64.StdEncoding.EncodeToString([]byte(customEncoded)), nil
}

func (e *CustomBase64Encryptor) Decrypt(ciphertext string) (string, error) {
	// 先解开外层标准Base64
	customEncoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 再用自定义字符表解码
	plaintext, err := e.encoding.DecodeString(string(customEncoded))
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
