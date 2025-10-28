package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/api"
	"github.com/nextkey/nextkey/backend/internal/crypto"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/pkg/config"
)

func main() {
	cfg := config.Load()

	if err := crypto.SetKey(cfg.Security.AESKey); err != nil {
		log.Fatalf("设置加密密钥失败: %v", err)
	}

	if err := database.Initialize(cfg.Database.Path); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	api.RegisterRoutes(router)

	distFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Printf("前端文件未找到，跳过静态文件服务: %v", err)
	} else {
		// 注册 /assets 路由用于静态资源
		assetsFS, err := fs.Sub(distFS, "assets")
		if err == nil {
			router.StaticFS("/assets", http.FS(assetsFS))
		}

		// SPA fallback: 所有未匹配路由返回 index.html
		router.NoRoute(func(c *gin.Context) {
			data, err := fs.ReadFile(distFS, "index.html")
			if err != nil {
				c.String(500, "无法读取前端文件")
				return
			}
			c.Data(200, "text/html; charset=utf-8", data)
		})
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("\nNextKey 服务启动成功\n")
	fmt.Printf("访问地址: http://localhost%s\n", addr)
	fmt.Printf("默认账号: admin / admin123\n\n")

	if err := router.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
