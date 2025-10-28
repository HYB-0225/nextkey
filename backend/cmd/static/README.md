# Static 静态资源目录

## 📁 目录说明

此目录用于存放前端构建产物，在构建时会被 Go embed 嵌入到二进制文件中。

## 🔄 使用流程

### 开发模式

后端和前端分离运行：

```bash
# 终端1: 启动后端
go run backend/cmd/main.go backend/cmd/embed.go

# 终端2: 启动前端开发服务器
cd frontend
npm run dev
```

此时访问:
- 后端 API: http://localhost:8080/api
- 前端界面: http://localhost:3000 (Vite 开发服务器)

### 生产构建

前端构建后嵌入到 Go 二进制：

```bash
# 使用构建脚本
./scripts/build.sh        # Linux/macOS
scripts\build.bat         # Windows

# 或手动构建
cd frontend
npm run build
cp -r dist/* ../backend/cmd/static/
cd ..
go build -o nextkey backend/cmd/main.go backend/cmd/embed.go
```

构建后此目录包含：
```
static/
├── index.html           # 前端入口
├── assets/              # JS/CSS/图片等资源
│   ├── index-xxx.js
│   ├── index-xxx.css
│   └── ...
└── .gitkeep            # Git 占位文件
```

## ⚠️ 注意事项

1. **不要提交构建产物到 Git**
   - `static/` 目录下除 `.gitkeep` 外的文件已在 `.gitignore` 中
   - 只提交源码，不提交构建产物

2. **embed 路径**
   - `backend/cmd/embed.go` 中的 `//go:embed static` 指向此目录
   - 路径相对于 `embed.go` 文件位置

3. **开发调试**
   - 开发时建议使用前端开发服务器，支持热重载
   - 只在测试完整构建或部署前构建前端

## 🚀 快速测试完整构建

```bash
# 1. 构建前端
cd frontend && npm run build && cd ..

# 2. 复制到 static
rm -rf backend/cmd/static/*
cp -r frontend/dist/* backend/cmd/static/

# 3. 运行后端（已嵌入前端）
go run backend/cmd/main.go backend/cmd/embed.go

# 4. 访问 http://localhost:8080
```

