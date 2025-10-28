package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/service"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

func GetCloudVar(c *gin.Context) {
	projectID, _ := c.Get("project_id")
	key := c.Param("key")

	cloudVarSvc := service.NewCloudVarService()
	cloudVar, err := cloudVarSvc.Get(projectID.(uint), key)
	if err != nil {
		utils.EncryptedError(c, 404, err.Error())
		return
	}

	utils.EncryptedSuccess(c, cloudVar)
}

func SetCloudVar(c *gin.Context) {
	var req service.CreateCloudVarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cloudVarSvc := service.NewCloudVarService()
	cloudVar, err := cloudVarSvc.Set(&req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, cloudVar)
}

func ListCloudVars(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.DefaultQuery("project_id", "0"))

	cloudVarSvc := service.NewCloudVarService()
	cloudVars, err := cloudVarSvc.List(uint(projectID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, cloudVars)
}

func DeleteCloudVar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	cloudVarSvc := service.NewCloudVarService()
	if err := cloudVarSvc.Delete(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func BatchSetCloudVars(c *gin.Context) {
	var req struct {
		Data []service.CreateCloudVarRequest `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cloudVarSvc := service.NewCloudVarService()
	if err := cloudVarSvc.BatchSet(req.Data); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "批量设置成功"})
}

func BatchDeleteCloudVars(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	cloudVarSvc := service.NewCloudVarService()
	if err := cloudVarSvc.BatchDelete(req.IDs); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "批量删除成功"})
}
