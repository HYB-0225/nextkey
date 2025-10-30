package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/crypto"
	"github.com/nextkey/nextkey/backend/internal/service"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

// GetEncryptionSchemes 获取支持的加密方案列表
func GetEncryptionSchemes(c *gin.Context) {
	schemes := crypto.ListSchemes()
	utils.Success(c, schemes)
}

// UpdateProjectEncryption 修改项目的加密方案
func UpdateProjectEncryption(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		EncryptionScheme string `json:"encryption_scheme"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if req.EncryptionScheme == "" {
		utils.Error(c, 400, "加密方案不能为空")
		return
	}

	projectSvc := service.NewProjectService()
	project, err := projectSvc.UpdateEncryptionScheme(uint(id), req.EncryptionScheme)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"encryption_scheme": project.EncryptionScheme,
		"encryption_key":    project.EncryptionKey,
	})
}

