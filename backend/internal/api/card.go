package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/middleware"
	"github.com/nextkey/nextkey/backend/internal/service"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

func UpdateCardCustomData(c *gin.Context) {
	cardID, exists := c.Get("card_id")
	if !exists {
		utils.EncryptedError(c, 400, "免费模式不支持此功能")
		return
	}

	var req struct {
		CustomData string `json:"custom_data"`
	}
	if err := middleware.GetDecryptedData(c, &req); err != nil {
		utils.EncryptedError(c, 400, "参数错误")
		return
	}

	cardSvc := service.NewCardService()
	if err := cardSvc.UpdateCustomData(cardID.(uint), req.CustomData); err != nil {
		utils.EncryptedError(c, 500, err.Error())
		return
	}

	utils.EncryptedSuccess(c, gin.H{"message": "更新成功"})
}

func CreateCards(c *gin.Context) {
	var req service.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cardSvc := service.NewCardService()
	cards, err := cardSvc.CreateBatch(&req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, cards)
}

func ListCards(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.DefaultQuery("project_id", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	cardSvc := service.NewCardService()
	cards, total, err := cardSvc.List(uint(projectID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":  cards,
		"total": total,
		"page":  page,
	})
}

func GetCard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	cardSvc := service.NewCardService()
	card, err := cardSvc.Get(uint(id))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	utils.Success(c, card)
}

func UpdateCard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req service.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cardSvc := service.NewCardService()
	card, err := cardSvc.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, card)
}

func DeleteCard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	cardSvc := service.NewCardService()
	if err := cardSvc.Delete(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func BatchUpdateCards(c *gin.Context) {
	var req struct {
		IDs  []uint                    `json:"ids"`
		Data service.UpdateCardRequest `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cardSvc := service.NewCardService()
	if err := cardSvc.BatchUpdate(req.IDs, &req.Data); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "批量更新成功"})
}

func BatchDeleteCards(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cardSvc := service.NewCardService()
	if err := cardSvc.BatchDelete(req.IDs); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "批量删除成功"})
}
