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
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path[1:]
			if path == "" {
				path = "index.html"
			}

			file, err := distFS.Open(path)
			if err != nil {
				indexFile, _ := distFS.Open("index.html")
				if indexFile != nil {
					c.FileFromFS("index.html", http.FS(distFS))
				}
				return
			}
			defer file.Close()
			c.FileFromFS(path, http.FS(distFS))
		})
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("\n🚀 NextKey 服务启动成功\n")
	fmt.Printf("📍 访问地址: http://localhost%s\n", addr)
	fmt.Printf("👤 默认账号: admin / admin123\n\n")

	if err := router.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
