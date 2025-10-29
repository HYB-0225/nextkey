# NextKey 部署指南

## 快速部署

### 1. 下载预构建版本

从 [Releases](https://github.com/HYB-0225/nextkey/releases) 下载对应平台的二进制文件。

### 2. 首次运行

```bash
# Linux/macOS
chmod +x nextkey
./nextkey

# Windows
nextkey.exe
```

首次运行会自动:
- 创建 `config.yaml` 配置文件
- 初始化 SQLite 数据库 `nextkey.db`
- 创建默认管理员账号: `admin / admin123`

### 3. 访问管理后台

浏览器访问: http://localhost:8080

## 配置说明

`config.yaml`:

```yaml
server:
  port: 8080              # 服务端口
  mode: release           # debug/release

database:
  path: ./nextkey.db      # 数据库路径

security:
  aes_key: "自动生成"     # AES加密密钥
  jwt_secret: "自动生成"  # JWT密钥
  token_expire: 3600      # Token有效期(秒)
  replay_window: 300      # 防重放时间窗口(秒)

admin:
  username: admin
  password: admin123      # 首次运行后请修改
```

## 生产环境部署

### 使用 systemd (Linux)

创建 `/etc/systemd/system/nextkey.service`:

```ini
[Unit]
Description=NextKey Service
After=network.target

[Service]
Type=simple
User=nextkey
WorkingDirectory=/opt/nextkey
ExecStart=/opt/nextkey/nextkey
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
sudo systemctl daemon-reload
sudo systemctl enable nextkey
sudo systemctl start nextkey
sudo systemctl status nextkey
```

### 使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### HTTPS 配置

使用 Let's Encrypt:

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## 数据备份

定期备份以下文件:
- `nextkey.db` - 数据库文件
- `config.yaml` - 配置文件

备份脚本示例:

```bash
#!/bin/bash
backup_dir="/backup/nextkey/$(date +%Y%m%d)"
mkdir -p $backup_dir
cp nextkey.db $backup_dir/
cp config.yaml $backup_dir/
```

## 性能优化

### SQLite 优化

对于高并发场景，考虑:
1. 定期执行 `VACUUM` 优化数据库
2. 启用 WAL 模式
3. 或迁移到 PostgreSQL/MySQL

### 限制请求频率

使用 Nginx 限流:

```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

location /api/ {
    limit_req zone=api burst=20;
    proxy_pass http://localhost:8080;
}
```

## 监控

### 日志

服务日志输出到标准输出，使用 systemd 查看:

```bash
sudo journalctl -u nextkey -f
```

### 健康检查

可以通过访问管理后台检查服务状态，或添加健康检查端点。

## 故障排查

### 端口被占用

```bash
# 查看端口占用
netstat -tlnp | grep 8080

# 修改配置文件中的端口
vim config.yaml
```

### 数据库锁定

SQLite 在高并发时可能出现锁定，考虑:
1. 减少并发请求
2. 使用连接池
3. 迁移到其他数据库

### 前端无法访问

检查:
1. 二进制文件是否包含前端资源
2. 浏览器控制台是否有错误
3. API 代理配置是否正确

## 安全建议

1. **修改默认密码**: 首次登录后立即修改管理员密码
2. **防火墙**: 仅开放必要端口
3. **HTTPS**: 生产环境必须使用 HTTPS
4. **定期更新**: 及时更新到最新版本
5. **备份**: 定期备份数据库
6. **密钥管理**: 妥善保管 AES 密钥和 JWT 密钥

