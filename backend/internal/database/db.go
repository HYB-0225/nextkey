package database

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/nextkey/nextkey/backend/internal/crypto"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/config"
	"github.com/nextkey/nextkey/backend/pkg/utils"
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

	configureSQLite()

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

	if err := DB.AutoMigrate(
		&models.Admin{},
		&models.AdminToken{},
		&models.AdminTokenBlacklist{},
		&models.Project{},
		&models.Card{},
		&models.Token{},
		&models.CloudVar{},
		&models.Nonce{},
		&models.UnbindRecord{},
	); err != nil {
		return err
	}

	// 迁移 Card 表缺失的列
	if err := migrateCardColumns(); err != nil {
		log.Printf("Card表迁移警告: %v", err)
	}

	// 迁移现有项目，为其生成加密密钥
	if err := migrateProjectEncryption(); err != nil {
		log.Printf("项目加密字段迁移警告: %v", err)
	}

	if err := migrateProjectUnbindSlug(); err != nil {
		log.Printf("项目解绑字段迁移警告: %v", err)
	}

	return nil
}

func configureSQLite() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		return
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := DB.Exec("PRAGMA busy_timeout = 5000").Error; err != nil {
		log.Printf("设置busy_timeout失败: %v", err)
	}

	var mode string
	if err := DB.Raw("PRAGMA journal_mode = WAL").Scan(&mode).Error; err != nil {
		log.Printf("设置WAL模式失败: %v", err)
	}
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

func migrateCardColumns() error {
	// 检查 cards 表是否存在
	if !DB.Migrator().HasTable(&models.Card{}) {
		return nil
	}

	// 定义需要检查的列及其默认值
	columns := map[string]string{
		"max_hwid":  "INTEGER DEFAULT -1",
		"max_ip":    "INTEGER DEFAULT -1",
		"hwid_list": "TEXT DEFAULT '[]'",
		"ip_list":   "TEXT DEFAULT '[]'",
		"frozen":    "INTEGER DEFAULT 0",
	}

	for columnName, columnDef := range columns {
		// 检查列是否存在
		var count int
		err := DB.Raw(`
			SELECT COUNT(*) 
			FROM pragma_table_info('cards') 
			WHERE name = ?
		`, columnName).Scan(&count).Error

		if err != nil {
			return err
		}

		// 如果列不存在，则添加
		if count == 0 {
			sql := "ALTER TABLE cards ADD COLUMN " + columnName + " " + columnDef
			if err := DB.Exec(sql).Error; err != nil {
				log.Printf("添加列 %s 失败: %v", columnName, err)
				return err
			}
			log.Printf("Card表迁移完成: 已添加列 %s", columnName)
		}
	}

	return nil
}

func migrateProjectEncryption() error {
	var projects []models.Project
	if err := DB.Find(&projects).Error; err != nil {
		return err
	}

	for _, p := range projects {
		if p.EncryptionKey == "" {
			p.EncryptionScheme = "aes-256-gcm"
			p.EncryptionKey = crypto.GenerateEncryptionKey()
			if err := DB.Save(&p).Error; err != nil {
				log.Printf("为项目 %s 生成加密密钥失败: %v", p.Name, err)
			} else {
				log.Printf("为项目 %s 生成加密密钥: %s", p.Name, p.EncryptionKey)
			}
		}
	}
	return nil
}

func migrateProjectUnbindSlug() error {
	if !DB.Migrator().HasTable(&models.Project{}) {
		return nil
	}

	if !DB.Migrator().HasColumn(&models.Project{}, "unbind_slug") {
		if err := DB.Exec("ALTER TABLE projects ADD COLUMN unbind_slug TEXT").Error; err != nil {
			return err
		}
	}

	var projects []models.Project
	if err := DB.Where("unbind_slug = '' OR unbind_slug IS NULL").Find(&projects).Error; err != nil {
		return err
	}

	for i := range projects {
		slug, err := generateUniqueUnbindSlug()
		if err != nil {
			return err
		}
		if err := DB.Model(&projects[i]).Update("unbind_slug", slug).Error; err != nil {
			return err
		}
	}

	return DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_unbind_slug ON projects(unbind_slug)").Error
}

func generateUniqueUnbindSlug() (string, error) {
	for i := 0; i < 5; i++ {
		slug := utils.RandomString(24, utils.CharsetTypeAlphanumeric)
		var count int64
		if err := DB.Model(&models.Project{}).Where("unbind_slug = ?", slug).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
	}
	return "", errors.New("生成解绑链接失败")
}

func cleanExpiredNonces() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		// 清理过期的Nonce(10分钟前)
		nonceCutoff := now.Add(-10 * time.Minute)
		DB.Where("created_at < ?", nonceCutoff).Delete(&models.Nonce{})

		// 清理过期的刷新令牌
		DB.Where("expire_at < ?", now).Delete(&models.AdminToken{})

		// 清理过期的JWT黑名单记录
		DB.Where("expire_at < ?", now).Delete(&models.AdminTokenBlacklist{})
	}
}
