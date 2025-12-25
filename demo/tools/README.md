# NextKey 测试工具

提供可视化图形界面的测试客户端，用于快速验证 NextKey API 对接。

## 功能特性

- ✅ **配置管理** - 保存/加载服务器配置（从 `config.yaml` 读取密钥仅兼容旧版本）
- ✅ **登录测试** - 测试卡密登录，支持HWID和IP参数
- ✅ **心跳验证** - 手动或自动（30秒间隔）心跳测试
- ✅ **云变量查询** - 实时查询云端变量值
- ✅ **专属信息更新** - 更新卡密专属JSON数据
- ✅ **项目信息** - 获取项目详细信息
- ✅ **请求查看器** - 查看加密前后的请求数据和响应详情
- ✅ **实时日志** - 所有操作实时日志记录，支持导出

## 环境要求

- Python 3.7+
- Tkinter (Python自带，无需额外安装)

## 安装

### 1. 克隆仓库（如果还没有）

```bash
git clone https://github.com/HYB-0225/nextkey.git
cd nextkey/tools
```

### 2. 安装依赖

```bash
pip install -r requirements.txt
```

或单独安装：

```bash
pip install pycryptodome requests pyyaml
```

## 使用方法

### 启动测试工具

```bash
python gui-test-client.py
```

或在 Windows 上双击运行：

```bash
# 先确保 .py 文件关联到 Python
gui-test-client.py
```

### 快速开始

#### 1. 配置服务器

切换到 **配置** 标签页：

1. 输入 **服务器URL**: `http://localhost:8080`
2. 输入 **项目UUID**: 从管理后台获取
3. 输入 **加密密钥**: 从管理后台项目详情获取（通常为64字符字符串）
4. 点击 **测试连接** 验证服务器可达性
5. 点击 **保存配置** 保存到本地文件

#### 2. 登录测试

切换到 **登录测试** 标签页：

1. 输入 **卡密**（必填）
2. 输入 **设备码** (可选，用于HWID验证)
3. 输入 **IP地址** (可选，留空则使用请求IP)
4. 点击 **登录**
5. 查看 Token 信息和响应详情

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGc...",
    "expire_at": "2024-01-01T12:00:00Z",
    "card": {
      "id": 1,
      "card_key": "TEST-123456",
      "activated": true,
      "duration": 2592000,
      "custom_data": "{}"
    }
  }
}
```

#### 3. API 测试

切换到 **API测试** 标签页：

**心跳测试**:
- 点击 **手动心跳** 发送一次心跳
- 点击 **开始自动心跳** 每30秒自动发送
- 点击 **停止自动心跳** 停止自动心跳

**云变量查询**:
1. 输入变量 Key (如: `notice`)
2. 点击 **查询**
3. 查看返回的变量值

**专属信息更新**:
1. 在文本框中编辑 JSON 数据，例如:
   ```json
   {
     "user_level": 5,
     "points": 1000,
     "last_login": "2024-01-01"
   }
   ```
2. 点击 **更新**

**项目信息查询**:
- 点击 **获取项目信息** 查看项目详情

#### 4. 查看日志

切换到 **日志** 标签页：

- 实时查看所有操作日志
- 成功操作显示为 <span style="color:green">绿色</span>
- 失败操作显示为 <span style="color:red">红色</span>
- 信息提示显示为 <span style="color:blue">蓝色</span>
- 点击 **导出日志** 保存到文件

## 配置文件

工具会自动创建 `nextkey_client_config.json` 保存配置：

```json
{
  "server_url": "http://localhost:8080",
  "project_uuid": "your-project-uuid",
  "aes_key": "your-aes-key"
}
```

**注意**: 此文件包含敏感密钥，请勿提交到版本控制系统。

## 从 config.yaml 读取密钥（已弃用）

当前版本的密钥为**项目级配置**，不再存放在 `config.yaml` 中。
旧版本若存在 `security.aes_key` 字段，可继续使用该按钮读取，作为兼容手段。

## 常见问题

### 1. 导入错误: No module named 'Crypto'

**解决方案**:
```bash
pip install pycryptodome
```

**注意**: 使用 `pycryptodome`，不是 `pycrypto`（已废弃）。

### 2. 连接失败

**检查清单**:
- [ ] 服务器是否已启动？
- [ ] URL 是否正确（包括 http:// 前缀）？
- [ ] 端口是否正确（默认8080）？
- [ ] 防火墙是否放行？

### 3. 登录失败: "解密失败"

**原因**: AES密钥不正确

**解决方案**:
1. 确认使用的是管理后台项目详情中的加密密钥
2. 密钥长度应为 32 或 64 字符
3. 旧版本可尝试从 `config.yaml` 读取（仅兼容旧配置）

### 4. 登录失败: "timestamp expired"

**原因**: 客户端与服务器时间差超过5分钟

**解决方案**:
```bash
# Windows
w32tm /resync

# Linux
sudo ntpdate time.windows.com
```

### 5. 心跳失败: "未授权"

**原因**: Token 已过期或无效

**解决方案**:
1. 重新登录获取新 Token
2. 检查 Token 有效期设置

## 高级功能

### 批量测试

可以编写脚本调用 `NextKeyClient` 类进行批量测试：

```python
from gui_test_client import NextKeyClient

client = NextKeyClient(
    "http://localhost:8080",
    "project-uuid",
    "aes-key"
)

# 批量测试卡密
card_keys = ["CARD1", "CARD2", "CARD3"]
for card in card_keys:
    try:
        result, _ = client.login(card)
        print(f"{card}: {'成功' if result['code'] == 0 else '失败'}")
    except Exception as e:
        print(f"{card}: 异常 - {e}")
```

### 自定义配置文件位置

修改 `NextKeyGUI.__init__()` 中的 `self.config_file`:

```python
self.config_file = "custom_config.json"  # 自定义路径
```

## 技术细节

### AES-256-GCM 加密

工具使用与服务端相同的加密方式：

```python
from Crypto.Cipher import AES
import base64

cipher = AES.new(key, AES.MODE_GCM)
ciphertext, tag = cipher.encrypt_and_digest(plaintext.encode())
encrypted = cipher.nonce + ciphertext + tag
result = base64.b64encode(encrypted).decode()
```

### 请求格式

外层请求体：

```json
{
  "timestamp": 1698505200,
  "nonce": "随机32字符串",
  "data": "Base64编码的加密数据"
}
```

`data` 解密后的内部结构：

```json
{
  "nonce": "同外层nonce",
  "timestamp": 1698505200,
  "data": {
    "...": "业务参数"
  }
}
```

## 开发贡献

欢迎提交改进建议和 Bug 报告！

**改进方向**:
- [ ] 支持批量导入卡密测试
- [ ] 添加性能测试模块
- [ ] 支持多项目切换
- [ ] 添加请求重放功能

## 相关文档

- [客户端对接文档](../docs/CLIENT.md)
- [API文档](../docs/API.md)
- [部署指南](../docs/DEPLOY.md)

## License

MIT License

