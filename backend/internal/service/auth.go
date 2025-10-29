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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
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

	// 生成JTI和刷新令牌
	jti := uuid.New().String()
	refreshToken := uuid.New().String()

	// 访问令牌有效期15分钟
	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	// 刷新令牌有效期7天
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)

	// 生成访问令牌(JWT)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"jti":      jti,
		"exp":      accessTokenExpiry.Unix(),
	})

	accessTokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// 保存刷新令牌到数据库
	adminToken := &models.AdminToken{
		AdminID:      admin.ID,
		RefreshToken: refreshToken,
		JTI:          jti,
		ExpireAt:     refreshTokenExpiry,
	}

	if err := database.DB.Create(adminToken).Error; err != nil {
		return nil, err
	}

	return &AdminLoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15分钟 = 900秒
	}, nil
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *AuthService) RefreshToken(req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 查找刷新令牌
	var adminToken models.AdminToken
	if err := database.DB.Where("refresh_token = ?", req.RefreshToken).Preload("Admin").First(&adminToken).Error; err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 检查刷新令牌是否过期
	if adminToken.IsExpired() {
		database.DB.Delete(&adminToken)
		return nil, errors.New("刷新令牌已过期")
	}

	// 验证管理员是否仍然存在
	if adminToken.Admin == nil {
		database.DB.Delete(&adminToken)
		return nil, errors.New("管理员不存在")
	}

	// 生成新的JTI和刷新令牌(令牌轮换)
	newJTI := uuid.New().String()
	newRefreshToken := uuid.New().String()

	// 访问令牌有效期15分钟
	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	// 刷新令牌有效期7天
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)

	// 生成新的访问令牌(JWT)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": adminToken.AdminID,
		"username": adminToken.Admin.Username,
		"jti":      newJTI,
		"exp":      accessTokenExpiry.Unix(),
	})

	accessTokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// 删除旧的刷新令牌
	database.DB.Delete(&adminToken)

	// 保存新的刷新令牌到数据库
	newAdminToken := &models.AdminToken{
		AdminID:      adminToken.AdminID,
		RefreshToken: newRefreshToken,
		JTI:          newJTI,
		ExpireAt:     refreshTokenExpiry,
	}

	if err := database.DB.Create(newAdminToken).Error; err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: newRefreshToken,
		ExpiresIn:    900, // 15分钟 = 900秒
	}, nil
}

type LogoutRequest struct {
	JTI string `json:"jti"`
}

func (s *AuthService) Logout(adminID uint, jti string) error {
	// 删除该管理员的所有刷新令牌(全局注销)
	if err := database.DB.Where("admin_id = ?", adminID).Delete(&models.AdminToken{}).Error; err != nil {
		return err
	}

	// 将当前访问令牌的JTI加入黑名单(15分钟后自动清理)
	blacklist := &models.AdminTokenBlacklist{
		JTI:      jti,
		ExpireAt: time.Now().Add(15 * time.Minute),
	}

	if err := database.DB.Create(blacklist).Error; err != nil {
		return err
	}

	return nil
}
