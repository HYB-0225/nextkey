package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
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
}

func (s *ProjectService) Create(req *CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		UUID:             uuid.New().String(),
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
	}

	if err := database.DB.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) List() ([]models.Project, error) {
	var projects []models.Project
	if err := database.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *ProjectService) GetByUUID(uuid string) (*models.Project, error) {
	var project models.Project
	if err := database.DB.Where("uuid = ?", uuid).First(&project).Error; err != nil {
		return nil, errors.New("项目不存在")
	}
	return &project, nil
}

func (s *ProjectService) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return nil, errors.New("项目不存在")
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
		project := &models.Project{
			UUID:             uuid.New().String(),
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
