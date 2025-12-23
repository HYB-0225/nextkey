package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/auth/login", middleware.LoginRateLimitMiddleware(), middleware.DecryptMiddleware(), CardLogin)
		api.GET("/crypto/schemes", GetEncryptionSchemes)
		api.POST("/card/unbind", middleware.DecryptMiddleware(), UnbindCardHWID)
		api.POST("/card/unbind-public", UnbindCardHWIDPublic)

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 需要加密的请求
			authenticated.POST("/heartbeat", middleware.DecryptMiddleware(), Heartbeat)
			authenticated.POST("/card/custom-data", middleware.DecryptMiddleware(), UpdateCardCustomData)
			authenticated.GET("/cloud-var/:key", middleware.DecryptMiddleware(), GetCloudVar)
			authenticated.POST("/cloud-var/:key", middleware.DecryptMiddleware(), GetCloudVar)
			authenticated.GET("/project/info", middleware.DecryptMiddleware(), GetProjectInfo)
			authenticated.POST("/project/info", middleware.DecryptMiddleware(), GetProjectInfo)
		}
	}

	admin := r.Group("/admin")
	{
		admin.POST("/login", middleware.LoginRateLimitMiddleware(), AdminLogin)
		admin.POST("/refresh", AdminRefreshToken)

		adminAuth := admin.Group("")
		adminAuth.Use(middleware.AdminAuthMiddleware())
		{
			adminAuth.POST("/logout", AdminLogout)

			adminAuth.GET("/projects", ListProjects)
			adminAuth.POST("/projects", CreateProject)
			adminAuth.PUT("/projects/:id", UpdateProject)
			adminAuth.DELETE("/projects/:id", DeleteProject)
			adminAuth.GET("/projects/:uuid", GetProjectByUUID)
			adminAuth.POST("/projects/batch", BatchCreateProjects)
			adminAuth.DELETE("/projects/batch", BatchDeleteProjects)
			adminAuth.POST("/projects/:id/encryption", UpdateProjectEncryption)

			adminAuth.GET("/cards", ListCards)
			adminAuth.POST("/cards", CreateCards)
			adminAuth.GET("/cards/:id", GetCard)
			adminAuth.PUT("/cards/:id", UpdateCard)
			adminAuth.DELETE("/cards/:id", DeleteCard)
			adminAuth.PUT("/cards/:id/freeze", FreezeCard)
			adminAuth.PUT("/cards/:id/unfreeze", UnfreezeCard)
			adminAuth.PUT("/cards/batch", BatchUpdateCards)
			adminAuth.DELETE("/cards/batch", BatchDeleteCards)
			adminAuth.PUT("/cards/batch/freeze", BatchFreezeCards)
			adminAuth.PUT("/cards/batch/unfreeze", BatchUnfreezeCards)

			adminAuth.GET("/cloud-vars", ListCloudVars)
			adminAuth.POST("/cloud-vars", SetCloudVar)
			adminAuth.DELETE("/cloud-vars/:id", DeleteCloudVar)
			adminAuth.POST("/cloud-vars/batch", BatchSetCloudVars)
			adminAuth.DELETE("/cloud-vars/batch", BatchDeleteCloudVars)
		}
	}
}
