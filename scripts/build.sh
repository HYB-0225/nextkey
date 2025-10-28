#!/bin/bash

set -e

echo "🚀 开始构建 NextKey..."

echo "📦 1. 构建前端..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi
npm run build
cd ..

echo "📋 2. 复制前端到静态目录..."
rm -rf backend/cmd/static/*
cp -r frontend/dist/* backend/cmd/static/

echo "🔨 3. 构建后端..."

VERSION=${VERSION:-"1.0.0"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS="-s -w"

mkdir -p build

echo "构建 Windows amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o build/nextkey-windows-amd64.exe backend/cmd/main.go backend/cmd/embed.go

echo "构建 Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o build/nextkey-linux-amd64 backend/cmd/main.go backend/cmd/embed.go

echo "构建 macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o build/nextkey-darwin-amd64 backend/cmd/main.go backend/cmd/embed.go

echo "构建 macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -o build/nextkey-darwin-arm64 backend/cmd/main.go backend/cmd/embed.go

echo "✅ 构建完成!"
echo "📁 二进制文件位于 build/ 目录"
ls -lh build/

