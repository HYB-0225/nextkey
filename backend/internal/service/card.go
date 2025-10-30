package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

// escapeLikeString 转义LIKE查询中的特殊字符
func escapeLikeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

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
	Duration   *int                `json:"duration"`
	Note       *string             `json:"note"`
	CardType   *string             `json:"card_type"`
	MaxHWID    *int                `json:"max_hwid"`
	MaxIP      *int                `json:"max_ip"`
	CustomData *string             `json:"custom_data"`
	HWIDList   *models.StringArray `json:"hwid_list"`
	IPList     *models.StringArray `json:"ip_list"`
}

type CardListFilter struct {
	ProjectID  uint
	Keyword    string
	CardType   string
	Note       string
	CustomData string
	Activated  string
	Frozen     string
	HWID       string
	IP         string
	StartTime  string
	EndTime    string
	Page       int
	PageSize   int
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

func (s *CardService) ListWithFilter(filter *CardListFilter) ([]models.Card, int64, error) {
	var cards []models.Card
	var total int64

	query := database.DB.Model(&models.Card{})

	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}

	if filter.Keyword != "" {
		// 转义LIKE特殊字符以防止SQL注入
		escapedKeyword := escapeLikeString(filter.Keyword)
		query = query.Where("card_key LIKE ? ESCAPE '\\'", "%"+escapedKeyword+"%")
	}

	if filter.CardType != "" {
		query = query.Where("card_type = ?", filter.CardType)
	}

	if filter.Note != "" {
		escapedNote := escapeLikeString(filter.Note)
		query = query.Where("note LIKE ? ESCAPE '\\'", "%"+escapedNote+"%")
	}

	if filter.CustomData != "" {
		escapedCustomData := escapeLikeString(filter.CustomData)
		query = query.Where("custom_data LIKE ? ESCAPE '\\'", "%"+escapedCustomData+"%")
	}

	if filter.HWID != "" {
		escapedHWID := escapeLikeString(filter.HWID)
		query = query.Where("hwid_list LIKE ? ESCAPE '\\'", "%"+escapedHWID+"%")
	}

	if filter.IP != "" {
		escapedIP := escapeLikeString(filter.IP)
		query = query.Where("ip_list LIKE ? ESCAPE '\\'", "%"+escapedIP+"%")
	}

	if filter.Activated == "true" {
		query = query.Where("activated = ?", true)
	} else if filter.Activated == "false" {
		query = query.Where("activated = ?", false)
	}

	if filter.Frozen == "true" {
		query = query.Where("frozen = ?", true)
	} else if filter.Frozen == "false" {
		query = query.Where("frozen = ?", false)
	}

	if filter.StartTime != "" {
		query = query.Where("created_at >= ?", filter.StartTime)
	}

	if filter.EndTime != "" {
		query = query.Where("created_at <= ?", filter.EndTime)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&cards).Error; err != nil {
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
		// 如果卡密已激活，从激活时间重新计算到期时间
		if card.Activated && card.ActivatedAt != nil {
			if card.Duration > 0 {
				expireAt := card.ActivatedAt.Add(time.Duration(card.Duration) * time.Second)
				card.ExpireAt = &expireAt
			} else {
				card.ExpireAt = nil
			}
		}
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
	if req.HWIDList != nil {
		card.HWIDList = *req.HWIDList
	}
	if req.IPList != nil {
		card.IPList = *req.IPList
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

	// 如果修改了 Duration，需要单独处理每个卡密以重新计算 ExpireAt
	if req.Duration != nil {
		var cards []models.Card
		if err := tx.Where("id IN ?", ids).Find(&cards).Error; err != nil {
			tx.Rollback()
			return err
		}

		for i := range cards {
			cards[i].Duration = *req.Duration
			// 如果卡密已激活，从激活时间重新计算到期时间
			if cards[i].Activated && cards[i].ActivatedAt != nil {
				if cards[i].Duration > 0 {
					expireAt := cards[i].ActivatedAt.Add(time.Duration(cards[i].Duration) * time.Second)
					cards[i].ExpireAt = &expireAt
				} else {
					cards[i].ExpireAt = nil
				}
			}

			if req.Note != nil {
				cards[i].Note = *req.Note
			}
			if req.CardType != nil {
				cards[i].CardType = *req.CardType
			}
			if req.MaxHWID != nil {
				cards[i].MaxHWID = *req.MaxHWID
			}
			if req.MaxIP != nil {
				cards[i].MaxIP = *req.MaxIP
			}
			if req.CustomData != nil {
				cards[i].CustomData = *req.CustomData
			}
			if req.HWIDList != nil {
				cards[i].HWIDList = *req.HWIDList
			}
			if req.IPList != nil {
				cards[i].IPList = *req.IPList
			}

			if err := tx.Save(&cards[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		// 如果没有修改 Duration，可以使用批量更新
		updates := make(map[string]interface{})
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
		if req.HWIDList != nil {
			updates["hwid_list"] = *req.HWIDList
		}
		if req.IPList != nil {
			updates["ip_list"] = *req.IPList
		}

		if len(updates) == 0 {
			tx.Rollback()
			return errors.New("没有需要更新的字段")
		}

		if err := tx.Model(&models.Card{}).Where("id IN ?", ids).Updates(updates).Error; err != nil {
			tx.Rollback()
			return err
		}
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

func (s *CardService) FreezeCard(id uint) error {
	var card models.Card
	if err := database.DB.First(&card, id).Error; err != nil {
		return errors.New("卡密不存在")
	}

	if card.Frozen {
		return errors.New("卡密已处于冻结状态")
	}

	card.Frozen = true
	return database.DB.Save(&card).Error
}

func (s *CardService) UnfreezeCard(id uint) error {
	var card models.Card
	if err := database.DB.First(&card, id).Error; err != nil {
		return errors.New("卡密不存在")
	}

	if !card.Frozen {
		return errors.New("卡密未被冻结")
	}

	card.Frozen = false
	return database.DB.Save(&card).Error
}

func (s *CardService) BatchFreeze(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("未选择卡密")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Card{}).Where("id IN ?", ids).Update("frozen", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *CardService) BatchUnfreeze(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("未选择卡密")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Card{}).Where("id IN ?", ids).Update("frozen", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

type UnbindRequest struct {
	CardKey     string `json:"card_key"`
	HWID        string `json:"hwid"`
	ProjectUUID string `json:"project_uuid"`
}

func (s *CardService) UnbindHWID(req *UnbindRequest) error {
	var project models.Project
	if err := database.DB.Where("uuid = ?", req.ProjectUUID).First(&project).Error; err != nil {
		return errors.New("项目不存在")
	}

	if !project.EnableUnbind {
		return errors.New("该项目未启用解绑功能")
	}

	var card models.Card
	if err := database.DB.Where("card_key = ? AND project_id = ?", req.CardKey, project.ID).First(&card).Error; err != nil {
		return errors.New("卡密不存在")
	}

	if card.IsFrozen() {
		return errors.New("卡密已冻结")
	}

	if req.HWID == "" {
		return errors.New("HWID不能为空")
	}

	var lastUnbind models.UnbindRecord
	if err := database.DB.Where("card_id = ?", card.ID).Order("unbind_at DESC").First(&lastUnbind).Error; err == nil {
		cooldownEnd := lastUnbind.UnbindAt.Add(time.Duration(project.UnbindCooldown) * time.Second)
		if time.Now().Before(cooldownEnd) {
			remainingTime := int(time.Until(cooldownEnd).Seconds())
			return errors.New("解绑冷却中，请等待 " + strconv.Itoa(remainingTime) + " 秒后再试")
		}
	}

	if project.UnbindVerifyHWID {
		found := false
		for _, hwid := range card.HWIDList {
			if hwid == req.HWID {
				found = true
				break
			}
		}
		if !found {
			return errors.New("该设备未绑定到此卡密")
		}
	}

	newHWIDList := make(models.StringArray, 0)
	for _, hwid := range card.HWIDList {
		if hwid != req.HWID {
			newHWIDList = append(newHWIDList, hwid)
		}
	}
	card.HWIDList = newHWIDList

	if project.UnbindDeductTime > 0 && card.ExpireAt != nil {
		newExpireAt := card.ExpireAt.Add(-time.Duration(project.UnbindDeductTime) * time.Second)
		if newExpireAt.Before(time.Now()) {
			card.ExpireAt = nil
		} else {
			card.ExpireAt = &newExpireAt
		}
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(&card).Error; err != nil {
		tx.Rollback()
		return err
	}

	unbindRecord := models.UnbindRecord{
		CardID:       card.ID,
		HWID:         req.HWID,
		UnbindAt:     time.Now(),
		DeductedTime: project.UnbindDeductTime,
	}
	if err := tx.Create(&unbindRecord).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
