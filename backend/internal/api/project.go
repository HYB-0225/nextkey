package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/service"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

func GetProjectInfo(c *gin.Context) {
	projectID, _ := c.Get("project_id")

	projectSvc := service.NewProjectService()
	project, err := projectSvc.GetByID(projectID.(uint))
	if err != nil {
		utils.EncryptedError(c, 404, err.Error())
		return
	}

	utils.EncryptedSuccess(c, gin.H{
		"uuid":       project.UUID,
		"name":       project.Name,
		"version":    project.Version,
		"update_url": project.UpdateURL,
	})
}

func CreateProject(c *gin.Context) {
	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if req.TokenExpire == 0 {
		req.TokenExpire = 3600
	}

	projectSvc := service.NewProjectService()
	project, err := projectSvc.Create(&req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, project)
}

func ListProjects(c *gin.Context) {
	projectSvc := service.NewProjectService()
	projects, err := projectSvc.List()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, projects)
}

func GetProjectByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	projectSvc := service.NewProjectService()
	project, err := projectSvc.GetByUUID(uuid)
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	utils.Success(c, project)
}

func UpdateProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	projectSvc := service.NewProjectService()
	project, err := projectSvc.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, project)
}

func DeleteProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	projectSvc := service.NewProjectService()
	if err := projectSvc.Delete(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func BatchCreateProjects(c *gin.Context) {
	var req struct {
		Data []service.CreateProjectRequest `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	projectSvc := service.NewProjectService()
	projects, err := projectSvc.BatchCreate(req.Data)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, projects)
}

func BatchDeleteProjects(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	projectSvc := service.NewProjectService()
	if err := projectSvc.BatchDelete(req.IDs); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "批量删除成功"})
}
