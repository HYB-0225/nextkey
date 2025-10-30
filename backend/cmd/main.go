package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/api"
	_ "github.com/nextkey/nextkey/backend/internal/crypto"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/middleware"
	"github.com/nextkey/nextkey/backend/internal/service"
	"github.com/nextkey/nextkey/backend/pkg/config"
)

func main() {
	// 检查配置文件权限
	checkConfigPermissions()

	cfg := config.Load()

	// 设置JWT密钥
	middleware.SetJWTSecret(cfg.Security.JWTSecret)
	service.SetJWTSecret(cfg.Security.JWTSecret)

	if err := database.Initialize(cfg.Database.Path, cfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 添加CORS中间件
	router.Use(middleware.CORSMiddleware())

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

		// SPA fallback: 只对请求 HTML 的请求返回 index.html
		router.NoRoute(func(c *gin.Context) {
			// 只在浏览器请求页面时返回 index.html
			if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
				data, err := fs.ReadFile(distFS, "index.html")
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
				c.Data(http.StatusOK, "text/html; charset=utf-8", data)
				return
			}
			// 其他资源请求返回 404
			c.Status(http.StatusNotFound)
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

func checkConfigPermissions() {
	configPath := "config.yaml"

	// 检查文件是否存在
	info, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// 文件不存在,将会生成默认配置
		return
	}
	if err != nil {
		log.Printf("警告: 无法检查配置文件权限: %v", err)
		return
	}

	// Windows系统跳过权限检查
	if runtime.GOOS == "windows" {
		return
	}

	// Unix/Linux系统检查文件权限
	mode := info.Mode()
	perm := mode.Perm()

	// 检查是否对组或其他用户可读
	if perm&0077 != 0 {
		log.Printf("警告: 配置文件 %s 权限过于宽松 (%s)", configPath, perm)
		log.Printf("建议运行: chmod 600 %s", configPath)
		log.Println("配置文件包含敏感密钥,应限制访问权限")
	}
}
