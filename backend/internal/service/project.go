package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nextkey/nextkey/backend/internal/crypto"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

type CreateProjectRequest struct {
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	EnableHWID       bool   `json:"enable_hwid"`
	EnableIP         bool   `json:"enable_ip"`
	Version          string `json:"version"`
	UpdateURL        string `json:"update_url"`
	TokenExpire      int    `json:"token_expire"`
	Description      string `json:"description"`
	EnableUnbind     bool   `json:"enable_unbind"`
	UnbindVerifyHWID bool   `json:"unbind_verify_hwid"`
	UnbindDeductTime int    `json:"unbind_deduct_time"`
	UnbindCooldown   int    `json:"unbind_cooldown"`
	EncryptionScheme string `json:"encryption_scheme"` // 加密方案，默认 aes-256-gcm
}

func (s *ProjectService) generateUnbindSlug() (string, error) {
	for i := 0; i < 5; i++ {
		slug := utils.RandomString(24, utils.CharsetTypeAlphanumeric)
		var count int64
		if err := database.DB.Model(&models.Project{}).Where("unbind_slug = ?", slug).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
	}
	return "", errors.New("生成解绑链接失败")
}

func (s *ProjectService) ensureUnbindSlug(project *models.Project) error {
	if project.UnbindSlug != "" {
		return nil
	}
	slug, err := s.generateUnbindSlug()
	if err != nil {
		return err
	}
	project.UnbindSlug = slug
	return database.DB.Model(project).Update("unbind_slug", slug).Error
}

func (s *ProjectService) Create(req *CreateProjectRequest) (*models.Project, error) {
	// 设置默认加密方案
	if req.EncryptionScheme == "" {
		req.EncryptionScheme = "aes-256-gcm"
	}

	// 验证加密方案是否支持
	if !crypto.SchemeExists(req.EncryptionScheme) {
		return nil, errors.New("不支持的加密方案: " + req.EncryptionScheme)
	}

	// 生成对应方案的密钥
	encryptionKey, err := crypto.GenerateKey(req.EncryptionScheme)
	if err != nil {
		return nil, err
	}

	unbindSlug, err := s.generateUnbindSlug()
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		UUID:             uuid.New().String(),
		UnbindSlug:       unbindSlug,
		Name:             req.Name,
		Mode:             req.Mode,
		EnableHWID:       req.EnableHWID,
		EnableIP:         req.EnableIP,
		Version:          req.Version,
		UpdateURL:        req.UpdateURL,
		TokenExpire:      req.TokenExpire,
		Description:      req.Description,
		EnableUnbind:     req.EnableUnbind,
		UnbindVerifyHWID: req.UnbindVerifyHWID,
		UnbindDeductTime: req.UnbindDeductTime,
		UnbindCooldown:   req.UnbindCooldown,
		EncryptionScheme: req.EncryptionScheme,
		EncryptionKey:    encryptionKey,
	}

	if err := database.DB.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) List(page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := database.DB.Model(&models.Project{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	for i := range projects {
		if err := s.ensureUnbindSlug(&projects[i]); err != nil {
			return nil, 0, err
		}
	}

	return projects, total, nil
}

func (s *ProjectService) GetByUUID(uuid string) (*models.Project, error) {
	var project models.Project
	if err := database.DB.Where("uuid = ?", uuid).First(&project).Error; err != nil {
		return nil, errors.New("项目不存在")
	}
	if err := s.ensureUnbindSlug(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return nil, errors.New("项目不存在")
	}
	if err := s.ensureUnbindSlug(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) Update(id uint, req *CreateProjectRequest) (*models.Project, error) {
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return nil, errors.New("项目不存在")
	}

	project.Name = req.Name
	project.Mode = req.Mode
	project.EnableHWID = req.EnableHWID
	project.EnableIP = req.EnableIP
	project.Version = req.Version
	project.UpdateURL = req.UpdateURL
	project.TokenExpire = req.TokenExpire
	project.Description = req.Description
	project.EnableUnbind = req.EnableUnbind
	project.UnbindVerifyHWID = req.UnbindVerifyHWID
	project.UnbindDeductTime = req.UnbindDeductTime
	project.UnbindCooldown = req.UnbindCooldown

	if err := database.DB.Save(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectService) Delete(id uint) error {
	return database.DB.Delete(&models.Project{}, id).Error
}

func (s *ProjectService) BatchCreate(reqs []CreateProjectRequest) ([]*models.Project, error) {
	if len(reqs) == 0 {
		return nil, errors.New("未提供项目数据")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	projects := make([]*models.Project, 0, len(reqs))

	for _, req := range reqs {
		// 设置默认加密方案
		if req.EncryptionScheme == "" {
			req.EncryptionScheme = "aes-256-gcm"
		}

		// 验证加密方案是否支持
		if !crypto.SchemeExists(req.EncryptionScheme) {
			tx.Rollback()
			return nil, errors.New("不支持的加密方案: " + req.EncryptionScheme)
		}

		// 生成对应方案的密钥
		encryptionKey, err := crypto.GenerateKey(req.EncryptionScheme)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		unbindSlug, err := s.generateUnbindSlug()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		project := &models.Project{
			UUID:             uuid.New().String(),
			UnbindSlug:       unbindSlug,
			Name:             req.Name,
			Mode:             req.Mode,
			EnableHWID:       req.EnableHWID,
			EnableIP:         req.EnableIP,
			Version:          req.Version,
			UpdateURL:        req.UpdateURL,
			TokenExpire:      req.TokenExpire,
			Description:      req.Description,
			EnableUnbind:     req.EnableUnbind,
			UnbindVerifyHWID: req.UnbindVerifyHWID,
			UnbindDeductTime: req.UnbindDeductTime,
			UnbindCooldown:   req.UnbindCooldown,
			EncryptionScheme: req.EncryptionScheme,
			EncryptionKey:    encryptionKey,
		}

		if err := tx.Create(project).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		projects = append(projects, project)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("未选择项目")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id IN ?", ids).Delete(&models.Project{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateEncryptionScheme 更新项目的加密方案
func (s *ProjectService) UpdateEncryptionScheme(id uint, scheme string) (*models.Project, error) {
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return nil, errors.New("项目不存在")
	}

	// 验证加密方案是否支持
	if !crypto.SchemeExists(scheme) {
		return nil, errors.New("不支持的加密方案: " + scheme)
	}

	// 生成新密钥
	encryptionKey, err := crypto.GenerateKey(scheme)
	if err != nil {
		return nil, err
	}

	// 更新项目
	project.EncryptionScheme = scheme
	project.EncryptionKey = encryptionKey

	if err := database.DB.Save(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectService) GetByUnbindSlug(slug string) (*models.Project, error) {
	var project models.Project
	if err := database.DB.Where("unbind_slug = ?", slug).First(&project).Error; err != nil {
		return nil, errors.New("项目不存在")
	}
	if err := s.ensureUnbindSlug(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

// GetOnlineCount 获取项目的在线人数
func (s *ProjectService) GetOnlineCount(projectID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Token{}).
		Where("project_id = ? AND expire_at > ? AND deleted_at IS NULL", projectID, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
