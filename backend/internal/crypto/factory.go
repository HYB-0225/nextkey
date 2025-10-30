package crypto

// GenerateEncryptionKey 生成默认加密方案的密钥（向后兼容）
// 推荐使用 GenerateKey(scheme) 代替
func GenerateEncryptionKey() string {
	key, _ := GenerateKey("aes-256-gcm")
	return key
}
