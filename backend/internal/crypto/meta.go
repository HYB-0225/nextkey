package crypto

// EncryptorMeta 加密方案元数据
type EncryptorMeta struct {
	Scheme        string `json:"scheme"`         // 加密方案标识
	Name          string `json:"name"`           // 显示名称
	Description   string `json:"description"`    // 描述信息
	SecurityLevel string `json:"security_level"` // 安全等级: secure/weak/insecure
	IsDeprecated  bool   `json:"is_deprecated"`  // 是否已废弃
}

// EncryptorFactory 加密器工厂
type EncryptorFactory struct {
	Meta         EncryptorMeta
	NewEncryptor func(key string) (Encryptor, error)
	GenerateKey  func() string
}

