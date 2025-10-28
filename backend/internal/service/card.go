package service

import (
	"errors"
	"time"

	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

type CardService struct{}

func NewCardService() *CardService {
	return &CardService{}
}

type CreateCardRequest struct {
	ProjectID uint   `json:"project_id"`
	CardKey   string `json:"card_key"`
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Count     int    `json:"count"`
	Duration  int    `json:"duration"`
	CardType  string `json:"card_type"`
	MaxHWID   int    `json:"max_hwid"`
	MaxIP     int    `json:"max_ip"`
	Note      string `json:"note"`
}

type UpdateCardRequest struct {
	Duration   *int    `json:"duration"`
	Note       *string `json:"note"`
	CardType   *string `json:"card_type"`
	MaxHWID    *int    `json:"max_hwid"`
	MaxIP      *int    `json:"max_ip"`
	CustomData *string `json:"custom_data"`
}

func (s *CardService) CreateBatch(req *CreateCardRequest) ([]models.Card, error) {
	var project models.Project
	if err := database.DB.First(&project, req.ProjectID).Error; err != nil {
		return nil, errors.New("项目不存在")
	}

	cards := make([]models.Card, 0, req.Count)

	for i := 0; i < req.Count; i++ {
		var cardKey string
		if req.CardKey != "" && req.Count == 1 {
			cardKey = req.CardKey
		} else {
			cardKey = utils.GenerateCardKey(req.Prefix, req.Suffix, 16)
		}

		card := models.Card{
			CardKey:   cardKey,
			ProjectID: req.ProjectID,
			Duration:  req.Duration,
			CardType:  req.CardType,
			MaxHWID:   req.MaxHWID,
			MaxIP:     req.MaxIP,
			Note:      req.Note,
			HWIDList:  make(models.StringArray, 0),
			IPList:    make(models.StringArray, 0),
		}

		if err := database.DB.Create(&card).Error; err != nil {
			continue
		}

		cards = append(cards, card)
	}

	if len(cards) == 0 {
		return nil, errors.New("创建卡密失败")
	}

	return cards, nil
}

func (s *CardService) List(projectID uint, page, pageSize int) ([]models.Card, int64, error) {
	var cards []models.Card
	var total int64

	query := database.DB.Model(&models.Card{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&cards).Error; err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}

func (s *CardService) Get(id uint) (*models.Card, error) {
	var card models.Card
	if err := database.DB.Preload("Project").First(&card, id).Error; err != nil {
		return nil, errors.New("卡密不存在")
	}
	return &card, nil
}

func (s *CardService) Update(id uint, req *UpdateCardRequest) (*models.Card, error) {
	var card models.Card
	if err := database.DB.First(&card, id).Error; err != nil {
		return nil, errors.New("卡密不存在")
	}

	if req.Duration != nil {
		card.Duration = *req.Duration
	}
	if req.Note != nil {
		card.Note = *req.Note
	}
	if req.CardType != nil {
		card.CardType = *req.CardType
	}
	if req.MaxHWID != nil {
		card.MaxHWID = *req.MaxHWID
	}
	if req.MaxIP != nil {
		card.MaxIP = *req.MaxIP
	}
	if req.CustomData != nil {
		card.CustomData = *req.CustomData
	}

	if err := database.DB.Save(&card).Error; err != nil {
		return nil, err
	}

	return &card, nil
}

func (s *CardService) Delete(id uint) error {
	return database.DB.Delete(&models.Card{}, id).Error
}

func (s *CardService) UpdateCustomData(cardID uint, customData string) error {
	var card models.Card
	if err := database.DB.First(&card, cardID).Error; err != nil {
		return errors.New("卡密不存在")
	}

	card.CustomData = customData
	return database.DB.Save(&card).Error
}

func (s *CardService) Heartbeat(cardID uint, projectID uint) error {
	var token models.Token
	if err := database.DB.Where("card_id = ? AND project_id = ?", cardID, projectID).
		Order("created_at DESC").First(&token).Error; err != nil {
		return errors.New("Token不存在")
	}

	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return errors.New("项目不存在")
	}

	token.ExpireAt = time.Now().Add(time.Duration(project.TokenExpire) * time.Second)
	return database.DB.Save(&token).Error
}

func (s *CardService) BatchUpdate(ids []uint, req *UpdateCardRequest) error {
	if len(ids) == 0 {
		return errors.New("未选择卡密")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.CardType != nil {
		updates["card_type"] = *req.CardType
	}
	if req.MaxHWID != nil {
		updates["max_hwid"] = *req.MaxHWID
	}
	if req.MaxIP != nil {
		updates["max_ip"] = *req.MaxIP
	}
	if req.CustomData != nil {
		updates["custom_data"] = *req.CustomData
	}

	if len(updates) == 0 {
		tx.Rollback()
		return errors.New("没有需要更新的字段")
	}

	if err := tx.Model(&models.Card{}).Where("id IN ?", ids).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *CardService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("未选择卡密")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id IN ?", ids).Delete(&models.Card{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
