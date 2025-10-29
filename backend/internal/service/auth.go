package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type LoginRequest struct {
	CardKey     string `json:"card_key"`
	HWID        string `json:"hwid,omitempty"`
	IP          string `json:"ip,omitempty"`
	ProjectUUID string `json:"project_uuid"`
}

type LoginResponse struct {
	Token    string       `json:"token"`
	ExpireAt time.Time    `json:"expire_at"`
	Card     *models.Card `json:"card"`
}

func (s *AuthService) CardLogin(req *LoginRequest) (*LoginResponse, error) {
	var project models.Project
	if err := database.DB.Where("uuid = ?", req.ProjectUUID).First(&project).Error; err != nil {
		return nil, errors.New("认证失败")
	}

	// 免费模式: 跳过所有验证,直接返回Token
	if project.Mode == "free" {
		tokenStr := uuid.New().String()
		expireAt := time.Now().Add(time.Duration(project.TokenExpire) * time.Second)

		token := &models.Token{
			Token:     tokenStr,
			CardID:    nil,
			ProjectID: project.ID,
			ExpireAt:  expireAt,
		}

		if err := database.DB.Create(token).Error; err != nil {
			return nil, err
		}

		return &LoginResponse{
			Token:    tokenStr,
			ExpireAt: expireAt,
			Card:     nil,
		}, nil
	}

	// 付费模式: 完整验证逻辑
	var card models.Card
	if err := database.DB.Where("card_key = ? AND project_id = ?", req.CardKey, project.ID).First(&card).Error; err != nil {
		return nil, errors.New("认证失败")
	}

	if !card.Activated {
		card.Activated = true
		activatedAt := time.Now()
		card.ActivatedAt = &activatedAt
		if card.Duration > 0 {
			expireAt := activatedAt.Add(time.Duration(card.Duration) * time.Second)
			card.ExpireAt = &expireAt
		}
		database.DB.Save(&card)
	}

	if card.IsExpired() {
		return nil, errors.New("认证失败")
	}

	if card.IsFrozen() {
		return nil, errors.New("认证失败")
	}

	// 验证设备码
	if project.EnableHWID {
		if req.HWID == "" {
			return nil, errors.New("认证失败")
		}
		found := false
		for _, hwid := range card.HWIDList {
			if hwid == req.HWID {
				found = true
				break
			}
		}
		if !found {
			if !card.CanAddHWID() {
				return nil, errors.New("认证失败")
			}
			card.HWIDList = append(card.HWIDList, req.HWID)
			database.DB.Save(&card)
		}
	}

	// 验证IP地址
	if project.EnableIP {
		if req.IP == "" {
			return nil, errors.New("认证失败")
		}
		found := false
		for _, ip := range card.IPList {
			if ip == req.IP {
				found = true
				break
			}
		}
		if !found {
			if !card.CanAddIP() {
				return nil, errors.New("认证失败")
			}
			card.IPList = append(card.IPList, req.IP)
			database.DB.Save(&card)
		}
	}

	tokenStr := uuid.New().String()
	expireAt := time.Now().Add(time.Duration(project.TokenExpire) * time.Second)

	cardID := card.ID
	token := &models.Token{
		Token:     tokenStr,
		CardID:    &cardID,
		ProjectID: project.ID,
		ExpireAt:  expireAt,
	}

	if err := database.DB.Create(token).Error; err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    tokenStr,
		ExpireAt: expireAt,
		Card:     &card,
	}, nil
}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func (s *AuthService) AdminLogin(req *AdminLoginRequest) (*AdminLoginResponse, error) {
	var admin models.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码 - 支持bcrypt和旧的SHA256格式(向后兼容)
	if !verifyPassword(admin.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 如果使用旧的SHA256格式,自动升级到bcrypt
	if !strings.HasPrefix(admin.Password, "$2a$") && !strings.HasPrefix(admin.Password, "$2b$") {
		newHash, err := hashPasswordBcrypt(req.Password)
		if err == nil {
			admin.Password = newHash
			database.DB.Save(&admin)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AdminLoginResponse{Token: tokenString}, nil
}

// hashPasswordBcrypt 使用bcrypt哈希密码
func hashPasswordBcrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// hashPasswordSHA256 使用SHA256哈希密码(仅用于向后兼容)
func hashPasswordSHA256(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// verifyPassword 验证密码,支持bcrypt和SHA256
func verifyPassword(hashedPassword, password string) bool {
	// 尝试bcrypt验证
	if strings.HasPrefix(hashedPassword, "$2a$") || strings.HasPrefix(hashedPassword, "$2b$") {
		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		return err == nil
	}
	// 向后兼容SHA256
	return hashedPassword == hashPasswordSHA256(password)
}
