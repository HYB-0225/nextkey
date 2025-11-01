package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/models"
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	projectSvc := service.NewProjectService()
	projects, total, err := projectSvc.List(page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// 为每个项目添加在线人数
	type ProjectWithOnline struct {
		*models.Project
		OnlineCount int64 `json:"online_count"`
	}

	projectsWithOnline := make([]ProjectWithOnline, len(projects))
	for i, project := range projects {
		onlineCount, _ := projectSvc.GetOnlineCount(project.ID)
		projectsWithOnline[i] = ProjectWithOnline{
			Project:     &projects[i],
			OnlineCount: onlineCount,
		}
	}

	utils.Success(c, gin.H{
		"list":  projectsWithOnline,
		"total": total,
		"page":  page,
	})
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
