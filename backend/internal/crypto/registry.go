package crypto

import (
	"errors"
	"sync"
)

var (
	registry = &encryptorRegistry{
		factories: make(map[string]*EncryptorFactory),
	}
)

type encryptorRegistry struct {
	mu        sync.RWMutex
	factories map[string]*EncryptorFactory
}

// Register 注册加密方案
func Register(factory EncryptorFactory) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.factories[factory.Meta.Scheme] = &factory
}

// NewEncryptor 创建加密器实例
func NewEncryptor(scheme, key string) (Encryptor, error) {
	registry.mu.RLock()
	factory, exists := registry.factories[scheme]
	registry.mu.RUnlock()

	if !exists {
		return nil, errors.New("unsupported encryption scheme: " + scheme)
	}

	return factory.NewEncryptor(key)
}

// GenerateKey 生成指定方案的密钥
func GenerateKey(scheme string) (string, error) {
	registry.mu.RLock()
	factory, exists := registry.factories[scheme]
	registry.mu.RUnlock()

	if !exists {
		return "", errors.New("unsupported encryption scheme: " + scheme)
	}

	return factory.GenerateKey(), nil
}

// ListSchemes 列出所有已注册的加密方案
func ListSchemes() []EncryptorMeta {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	schemes := make([]EncryptorMeta, 0, len(registry.factories))
	for _, factory := range registry.factories {
		schemes = append(schemes, factory.Meta)
	}
	return schemes
}

// GetMeta 获取指定方案的元数据
func GetMeta(scheme string) (*EncryptorMeta, error) {
	registry.mu.RLock()
	factory, exists := registry.factories[scheme]
	registry.mu.RUnlock()

	if !exists {
		return nil, errors.New("unsupported encryption scheme: " + scheme)
	}

	meta := factory.Meta
	return &meta, nil
}

// SchemeExists 检查加密方案是否存在
func SchemeExists(scheme string) bool {
	registry.mu.RLock()
	_, exists := registry.factories[scheme]
	registry.mu.RUnlock()
	return exists
}

