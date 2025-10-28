package database

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/nextkey/nextkey/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	if err := migrate(); err != nil {
		return err
	}

	if err := initializeDefaultAdmin(); err != nil {
		return err
	}

	go cleanExpiredNonces()

	return nil
}

func migrate() error {
	return DB.AutoMigrate(
		&models.Admin{},
		&models.Project{},
		&models.Card{},
		&models.Token{},
		&models.CloudVar{},
		&models.Nonce{},
	)
}

func initializeDefaultAdmin() error {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword := hashPassword("admin123")
	admin := &models.Admin{
		Username: "admin",
		Password: hashedPassword,
	}

	if err := DB.Create(admin).Error; err != nil {
		return err
	}

	log.Println("已创建默认管理员账号: admin / admin123")
	return nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func cleanExpiredNonces() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-10 * time.Minute)
		DB.Where("created_at < ?", cutoff).Delete(&models.Nonce{})
	}
}
