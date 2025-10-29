# Changelog

All notable changes to this project will be documented in this file.

## [v1.0.0] - 2025-01-29

### 🎉 首次发布

NextKey 卡密验证与云控制系统正式发布！

### ✨ 核心功能

- **卡密系统** - 支持自定义生成、设备码/IP验证、批量操作
- **高级搜索** - 多维度搜索卡密（关键词、类型、状态、设备码、IP等）
- **批量导出** - 支持JSON/TXT/CSV三种格式导出卡密
- **云变量** - 项目级别的云端变量存储
- **版本控制** - 客户端版本管理和更新控制
- **安全通信** - AES-256-GCM加密，防重放攻击
- **项目隔离** - 多项目管理，每个项目独立UUID
- **开箱即用** - 单一二进制文件，自动初始化
- **现代化UI** - 响应式设计，支持桌面端和移动端

### 🔧 技术特性

- **后端**: Go + Gin + SQLite
- **前端**: Vue 3 + Element Plus
- **加密**: AES-256-GCM + JWT
- **跨平台**: Linux/Windows/Android ARM64 支持

### 📦 发布包

- `nextkey_1.0.0_linux_amd64.tar.gz`
- `nextkey_1.0.0_linux_arm64.tar.gz`
- `nextkey_1.0.0_windows_amd64.zip`
- `nextkey_1.0.0_windows_arm64.zip`
- `nextkey_1.0.0_android_arm64.tar.gz`

### 📚 文档

- [README](README.md) - 项目介绍和快速开始
- [API文档](docs/API.md) - 详细的API接口说明
- [客户端对接文档](docs/CLIENT.md) - 完整的客户端接入指南
- [部署指南](docs/DEPLOY.md) - 生产环境部署指南

---

## 开发历史

### 主要提交记录

- 实现卡密冻结/解冻功能
- 支持编辑卡密的HWID和IP列表
- 添加高级卡密搜索、批量导出和创建卡密对话框
- 更新仓库克隆链接到新账号
- 优化代码性能和可读性
- 移除过时的OpenSpec相关文件

---

**完整更新日志**: https://github.com/HYB-0225/nextkey/commits/v1.0.0

