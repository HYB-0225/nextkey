# NextKey - 卡密验证与云控制系统

开源的卡密验证、云变量、版本控制服务，支持高度自定义化。

## 特性

- 🔐 卡密系统 - 支持自定义生成、设备码/IP验证、批量操作
- 🔍 高级搜索 - 多维度搜索卡密（关键词、类型、状态、设备码、IP等）
- 📤 批量导出 - 支持JSON/TXT/CSV三种格式导出卡密
- ☁️ 云变量 - 项目级别的云端变量存储
- 📦 版本控制 - 客户端版本管理和更新控制
- 🔒 安全通信 - AES-256-GCM加密，防重放攻击
- 🎯 项目隔离 - 多项目管理，每个项目独立UUID
- 🚀 开箱即用 - 单一二进制文件，自动初始化
- 💎 现代化UI - 响应式设计，支持桌面端和移动端

## 快速开始

### 下载预构建版本

从 [Releases](https://github.com/nextkey/nextkey/releases) 下载对应平台的二进制文件。

### 运行

**开发模式**:
```bash
# Windows
run.bat

# Linux/macOS
chmod +x run.sh
./run.sh

# 或手动运行
go run backend/cmd/main.go backend/cmd/embed.go
```

**生产模式** (使用预构建二进制):
```bash
# Windows
nextkey.exe

# Linux/macOS
chmod +x nextkey
./nextkey
```

首次运行会自动初始化数据库并创建默认管理员账号:
- 用户名: `admin`
- 密码: `admin123` (请立即修改)

访问管理后台: http://localhost:8080

## 从源码构建

### 环境要求

- Go 1.21+
- Node.js 18+
- npm/pnpm

### 构建步骤

```bash
# 1. 克隆仓库
git clone https://github.com/nextkey/nextkey.git
cd nextkey

# 2. 使用构建脚本(推荐)
./scripts/build.sh        # Linux/macOS
scripts\build.bat         # Windows

# 或手动构建
# 2a. 构建前端
cd frontend
npm install
npm run build
cd ..

# 2b. 复制前端到静态目录
cp -r frontend/dist/* backend/cmd/static/  # Linux/macOS
xcopy /E /I /Y frontend\dist\* backend\cmd\static\  # Windows

# 2c. 构建后端
go mod download
go build -o nextkey backend/cmd/main.go backend/cmd/embed.go

# 3. 运行
./nextkey
```

### 跨平台编译

```bash
# 使用构建脚本(自动处理前端构建和复制)
chmod +x scripts/build.sh
./scripts/build.sh        # Linux/macOS
scripts\build.bat         # Windows

# 使用goreleaser(自行配置android-ndk)
goreleaser build --snapshot --clean

# 或手动编译(需先构建前端并复制到 backend/cmd/static/)
# Windows
GOOS=windows GOARCH=amd64 go build -o nextkey-windows-amd64.exe backend/cmd/main.go backend/cmd/embed.go

# Linux
GOOS=linux GOARCH=amd64 go build -o nextkey-linux-amd64 backend/cmd/main.go backend/cmd/embed.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o nextkey-darwin-amd64 backend/cmd/main.go backend/cmd/embed.go
```

## 配置

默认配置文件 `config.yaml` (首次运行自动生成):

```yaml
server:
  port: 8080
  mode: release # debug/release

database:
  path: ./nextkey.db

security:
  aes_key: "auto-generated-32-chars-key"
  jwt_secret: "auto-generated-secret"
  token_expire: 3600 # 秒
  replay_window: 300 # 防重放时间窗口(秒)

admin:
  username: admin
  password: admin123
```

## 文档

- **[客户端对接文档](docs/CLIENT.md)** - 完整的客户端接入指南，包含密钥配置、加密流程、API调用、常见问题等
- **[API文档](docs/API.md)** - 详细的API接口说明
- **[部署指南](docs/DEPLOY.md)** - 生产环境部署指南
- **[客户端示例](docs/examples/)** - Python、Go等多语言示例代码

## 测试工具

提供图形化测试工具，快速验证API对接：

```bash
cd tools
pip install -r requirements.txt
python gui-test-client.py
```

详见 [工具使用说明](tools/README.md)

## 开发

```bash
# 后端开发
go run backend/cmd/main.go backend/cmd/embed.go

# 前端开发
cd frontend
npm run dev
```

## License

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

