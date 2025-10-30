package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Security SecurityConfig `yaml:"security"`
	Admin    AdminConfig    `yaml:"admin"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type SecurityConfig struct {
	JWTSecret    string `yaml:"jwt_secret"`
	TokenExpire  int    `yaml:"token_expire"`
	ReplayWindow int    `yaml:"replay_window"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Load() *Config {
	configPath := "config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := generateDefault()
		saveConfig(configPath, cfg)
		log.Printf("已生成默认配置文件: %s", configPath)
		return cfg
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	return &cfg
}

func generateDefault() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "release",
		},
		Database: DatabaseConfig{
			Path: "./nextkey.db",
		},
		Security: SecurityConfig{
			JWTSecret:    generateRandomKey(32),
			TokenExpire:  3600,
			ReplayWindow: 300,
		},
		Admin: AdminConfig{
			Username: "admin",
			Password: "admin123",
		},
	}
}

func generateRandomKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("生成随机密钥失败: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func saveConfig(path string, cfg *Config) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Fatalf("写入配置文件失败: %v", err)
	}
}
