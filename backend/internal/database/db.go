package database

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Initialize(dbPath string, cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbPath,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	if err := migrate(); err != nil {
		return err
	}

	if err := syncAdminFromConfig(cfg); err != nil {
		return err
	}

	go cleanExpiredNonces()

	return nil
}

func migrate() error {
	// 先尝试迁移 Token 表的 card_id 字段,使其允许为空
	if err := migrateTokenCardID(); err != nil {
		log.Printf("Token表迁移警告: %v", err)
	}

	return DB.AutoMigrate(
		&models.Admin{},
		&models.Project{},
		&models.Card{},
		&models.Token{},
		&models.CloudVar{},
		&models.Nonce{},
	)
}

func migrateTokenCardID() error {
	// 检查 tokens 表是否存在
	if !DB.Migrator().HasTable(&models.Token{}) {
		return nil
	}

	// 检查 card_id 字段是否允许为空
	var tableInfo struct {
		NotNull int
	}
	err := DB.Raw(`
		SELECT "notnull" as not_null 
		FROM pragma_table_info('tokens') 
		WHERE name = 'card_id'
	`).Scan(&tableInfo).Error

	if err != nil {
		return err
	}

	// 如果 card_id 已经允许为空,则跳过迁移
	if tableInfo.NotNull == 0 {
		return nil
	}

	// SQLite 不支持直接修改列约束,需要重建表
	// 创建临时表
	if err := DB.Exec(`
		CREATE TABLE tokens_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT UNIQUE NOT NULL,
			card_id INTEGER,
			project_id INTEGER NOT NULL,
			expire_at DATETIME NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error; err != nil {
		return err
	}

	// 复制数据
	if err := DB.Exec(`
		INSERT INTO tokens_new (id, token, card_id, project_id, expire_at, created_at, updated_at, deleted_at)
		SELECT id, token, card_id, project_id, expire_at, created_at, updated_at, deleted_at
		FROM tokens
	`).Error; err != nil {
		DB.Exec(`DROP TABLE IF EXISTS tokens_new`)
		return err
	}

	// 删除旧表
	if err := DB.Exec(`DROP TABLE tokens`).Error; err != nil {
		return err
	}

	// 重命名新表
	if err := DB.Exec(`ALTER TABLE tokens_new RENAME TO tokens`).Error; err != nil {
		return err
	}

	// 重建索引
	DB.Exec(`CREATE INDEX idx_tokens_card_id ON tokens(card_id)`)
	DB.Exec(`CREATE INDEX idx_tokens_project_id ON tokens(project_id)`)
	DB.Exec(`CREATE INDEX idx_tokens_deleted_at ON tokens(deleted_at)`)

	log.Println("Token表迁移完成: card_id字段已允许为空")
	return nil
}

func syncAdminFromConfig(cfg *config.Config) error {
	hashedPassword := hashPasswordBcrypt(cfg.Admin.Password)

	// 查找现有管理员账号
	var admin models.Admin
	err := DB.First(&admin).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在管理员账号，创建新的
		admin = models.Admin{
			Username: cfg.Admin.Username,
			Password: hashedPassword,
		}
		if err := DB.Create(&admin).Error; err != nil {
			return err
		}
		log.Printf("已创建管理员账号: %s (密码来自配置文件)", cfg.Admin.Username)
	} else if err != nil {
		return err
	} else {
		// 已存在管理员账号，同步更新为配置文件中的值
		admin.Username = cfg.Admin.Username

		// 只在密码格式不同时更新密码
		needsUpdate := false
		if strings.HasPrefix(admin.Password, "$2a$") || strings.HasPrefix(admin.Password, "$2b$") {
			// 当前是bcrypt,验证是否与配置密码匹配
			if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(cfg.Admin.Password)); err != nil {
				needsUpdate = true
			}
		} else {
			// 当前是SHA256,需要升级
			needsUpdate = true
		}

		if needsUpdate {
			admin.Password = hashedPassword
			log.Printf("已同步管理员账号: %s (密码已更新为配置文件中的值)", cfg.Admin.Username)
		} else {
			log.Printf("已同步管理员账号: %s", cfg.Admin.Username)
		}

		if err := DB.Save(&admin).Error; err != nil {
			return err
		}
	}

	return nil
}

func hashPasswordBcrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt哈希失败,回退到SHA256: %v", err)
		return hashPasswordSHA256(password)
	}
	return string(hash)
}

func hashPasswordSHA256(password string) string {
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
