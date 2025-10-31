package service

import (
	"errors"

	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
)

type CloudVarService struct{}

func NewCloudVarService() *CloudVarService {
	return &CloudVarService{}
}

type CreateCloudVarRequest struct {
	ProjectID uint   `json:"project_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (s *CloudVarService) Set(req *CreateCloudVarRequest) (*models.CloudVar, error) {
	var cloudVar models.CloudVar
	err := database.DB.Where("project_id = ? AND key = ?", req.ProjectID, req.Key).First(&cloudVar).Error

	if err != nil {
		cloudVar = models.CloudVar{
			ProjectID: req.ProjectID,
			Key:       req.Key,
			Value:     req.Value,
		}
		if err := database.DB.Create(&cloudVar).Error; err != nil {
			return nil, err
		}
	} else {
		cloudVar.Value = req.Value
		if err := database.DB.Save(&cloudVar).Error; err != nil {
			return nil, err
		}
	}

	return &cloudVar, nil
}

func (s *CloudVarService) Get(projectID uint, key string) (*models.CloudVar, error) {
	var cloudVar models.CloudVar
	if err := database.DB.Where("project_id = ? AND key = ?", projectID, key).First(&cloudVar).Error; err != nil {
		return nil, errors.New("变量不存在")
	}
	return &cloudVar, nil
}

func (s *CloudVarService) List(projectID uint, page, pageSize int) ([]models.CloudVar, int64, error) {
	var cloudVars []models.CloudVar
	var total int64

	query := database.DB.Model(&models.CloudVar{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&cloudVars).Error; err != nil {
		return nil, 0, err
	}

	return cloudVars, total, nil
}

func (s *CloudVarService) Delete(id uint) error {
	return database.DB.Delete(&models.CloudVar{}, id).Error
}

func (s *CloudVarService) BatchSet(reqs []CreateCloudVarRequest) error {
	if len(reqs) == 0 {
		return errors.New("未提供变量数据")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, req := range reqs {
		var cloudVar models.CloudVar
		err := tx.Where("project_id = ? AND key = ?", req.ProjectID, req.Key).First(&cloudVar).Error

		if err != nil {
			cloudVar = models.CloudVar{
				ProjectID: req.ProjectID,
				Key:       req.Key,
				Value:     req.Value,
			}
			if err := tx.Create(&cloudVar).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			cloudVar.Value = req.Value
			if err := tx.Save(&cloudVar).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *CloudVarService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("未选择变量")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id IN ?", ids).Delete(&models.CloudVar{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
