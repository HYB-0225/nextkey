@echo off
setlocal

echo 🚀 开始构建 NextKey...

echo 📦 1. 构建前端...
cd frontend
if not exist "node_modules" (
    echo 安装前端依赖...
    call npm install
)
call npm run build
cd ..

echo 📋 2. 复制前端到静态目录...
xcopy /E /I /Y frontend\dist\* backend\cmd\static\

echo 🔨 3. 构建后端...

if not exist "build" mkdir build

echo 构建 Windows amd64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o build\nextkey-windows-amd64.exe backend\cmd\main.go backend\cmd\embed.go

echo 构建 Linux amd64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o build\nextkey-linux-amd64 backend\cmd\main.go backend\cmd\embed.go

echo 构建 macOS amd64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o build\nextkey-darwin-amd64 backend\cmd\main.go backend\cmd\embed.go

echo 构建 macOS arm64...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o build\nextkey-darwin-arm64 backend\cmd\main.go backend\cmd\embed.go

echo ✅ 构建完成!
echo 📁 二进制文件位于 build\ 目录
dir build

endlocal

