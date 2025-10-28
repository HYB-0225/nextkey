# NextKey 客户端对接文档

完整的客户端接入指南，帮助开发者快速集成 NextKey 卡密验证系统。

## 目录

- [快速开始](#快速开始)
- [密钥配置](#密钥配置)
- [加密通信详解](#加密通信详解)
- [API调用流程](#api调用流程)
- [多语言示例](#多语言示例)
- [常见问题FAQ](#常见问题faq)

---

## 快速开始

### 5分钟快速对接流程

1. **获取配置信息**
   - 从服务端 `config.yaml` 获取 `aes_key`
   - 在管理后台创建项目，获取 `project_uuid`
   - 生成测试卡密

2. **安装依赖** (以Python为例)
   ```bash
   pip install pycryptodome requests
   ```

3. **编写客户端代码**
   ```python
   from client import NextKeyClient
   
   client = NextKeyClient(
       server_url="http://localhost:8080",
       project_uuid="your-project-uuid",
       aes_key="your-aes-key-from-config"
   )
   
   # 登录
   result = client.login("your-card-key", "your-hwid")
   
   # 心跳保活
   client.heartbeat()
   
   # 获取云变量
   value = client.get_cloud_var("key_name")
   ```

4. **运行测试**
   - 使用测试工具验证对接: `python tools/gui-test-client.py`

---

## 密钥配置

### AES密钥获取

AES加密密钥位于服务端 `config.yaml` 文件中：

```yaml
security:
  aes_key: 632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037
```

### 密钥格式说明

**服务端格式**: 64字符十六进制字符串（HEX编码）

**客户端使用**:
- **方式1**: 直接使用前32字节
  ```python
  aes_key = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037"
  key_bytes = aes_key[:32].encode()  # 取前32个字符转字节
  ```

- **方式2**: HEX解码为16字节 (不推荐，当前实现使用方式1)
  ```python
  import binascii
  aes_key = "632005a33ebb7619c1efd3853c7109f1..."
  key_bytes = binascii.unhexlify(aes_key)  # 32字节
  ```

**重要**: 根据服务端实现 (`backend/internal/crypto/aes.go:32`)，当密钥长度为64字符时，直接取前32字节作为AES密钥。

### 安全存储建议

❌ **不要做**:
- 硬编码密钥到源代码中
- 将密钥提交到版本控制系统
- 在日志中输出完整密钥

✅ **推荐做法**:
- 使用环境变量存储密钥
- 使用配置文件，并加入 `.gitignore`
- 生产环境使用密钥管理服务（KMS）
- 客户端可考虑混淆/加密存储密钥

```python
# 示例：从环境变量读取
import os
aes_key = os.environ.get('NEXTKEY_AES_KEY')
```

---

## 加密通信详解

客户端API（`/api/*` 路径）需要使用 **AES-256-GCM** 加密通信。

### 加密算法

- **算法**: AES-256-GCM (Galois/Counter Mode)
- **密钥长度**: 32字节 (256位)
- **Nonce**: 12字节随机值（由GCM自动生成）
- **编码**: Base64

### 请求格式

```json
{
  "timestamp": 1698505200,
  "nonce": "随机32字符串",
  "data": "Base64编码的AES加密数据"
}
```

**字段说明**:
- `timestamp`: Unix时间戳（秒），用于防止重放攻击
- `nonce`: 随机字符串（建议24-32字符），每次请求唯一
- `data`: 加密后的实际请求数据（Base64编码）

### 加密流程详解

#### Python实现

```python
from Crypto.Cipher import AES
import base64
import json
import time
import secrets

def encrypt_request(aes_key, request_data):
    """
    加密请求数据
    
    Args:
        aes_key: 32字节AES密钥
        request_data: 字典格式的请求数据
    
    Returns:
        完整的加密请求体
    """
    # 1. 将请求数据转为JSON字符串
    json_str = json.dumps(request_data)
    
    # 2. 创建AES-GCM加密器
    cipher = AES.new(aes_key, AES.MODE_GCM)
    
    # 3. 加密数据（同时生成认证标签）
    ciphertext, tag = cipher.encrypt_and_digest(json_str.encode())
    
    # 4. 组合: nonce + ciphertext + tag
    encrypted = cipher.nonce + ciphertext + tag
    
    # 5. Base64编码
    encrypted_b64 = base64.b64encode(encrypted).decode()
    
    # 6. 构造完整请求体
    return {
        "timestamp": int(time.time()),
        "nonce": secrets.token_urlsafe(24),  # 生成随机nonce
        "data": encrypted_b64
    }
```

#### Go实现

```go
func encrypt(aesKey []byte, plaintext string) (string, error) {
    // 1. 创建AES cipher
    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return "", err
    }
    
    // 2. 创建GCM模式
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    // 3. 生成随机nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    // 4. 加密（nonce会被自动添加到结果前面）
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    
    // 5. Base64编码
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

### 响应格式与Nonce验证

**重要安全机制**: 从v1.1开始，所有敏感接口（卡密登录、心跳、云变量、卡密专属信息）的响应都会回显客户端发送的Nonce，防止重放攻击。

#### 响应结构

```json
{
  "nonce": "客户端请求时发送的nonce",
  "data": "Base64编码的AES加密响应数据"
}
```

#### 响应验证流程

1. **记录发送的Nonce**: 客户端发送请求时，保存当前请求的 `nonce` 值
2. **接收响应**: 获取服务器返回的响应
3. **验证Nonce匹配**: 检查响应中的 `nonce` 字段是否与发送的一致
4. **Base64解码**: 解码响应中的 `data` 字段
5. **AES-GCM解密**: 使用密钥解密数据
6. **解析业务数据**: 解析得到标准响应格式（`code`、`message`、`data`）

#### Python示例

```python
def make_encrypted_request(self, endpoint, data):
    """发送加密请求并验证响应"""
    # 1. 生成并记住nonce
    request_nonce = secrets.token_urlsafe(24)
    
    # 2. 加密请求数据
    json_data = json.dumps(data)
    encrypted_data = self.encrypt(json_data)
    
    req_body = {
        "timestamp": int(time.time()),
        "nonce": request_nonce,  # 记住这个nonce
        "data": encrypted_data
    }
    
    # 3. 发送请求
    response = self.session.post(f"{self.server_url}{endpoint}", json=req_body)
    resp_json = response.json()
    
    # 4. 验证响应nonce
    if resp_json.get("nonce") != request_nonce:
        raise ValueError("响应Nonce不匹配，可能遭受重放攻击！")
    
    # 5. 解密响应数据
    decrypted = self.decrypt(resp_json["data"])
    
    # 6. 解析业务数据
    return json.loads(decrypted)
```

#### 安全性说明

- **防重放攻击**: 即使攻击者抓取了完整的响应数据包，也无法将其用于其他请求，因为每次请求的Nonce都不同
- **防静态注入**: 攻击者无法伪造合法的响应，因为必须知道客户端刚刚发送的Nonce
- **双向验证**: 结合请求Nonce验证（服务端）和响应Nonce验证（客户端），形成完整的防护链

### 解密流程

服务端处理流程：

1. **验证时间戳**: 检查 `timestamp` 是否在允许范围内（±5分钟）
2. **验证Nonce**: 检查 `nonce` 是否在replay_window内已使用过
3. **Base64解码**: 解码 `data` 字段
4. **提取组件**: 分离 nonce(12字节) + tag(16字节) + ciphertext
5. **AES解密**: 使用GCM模式解密并验证
6. **解析JSON**: 将解密后的数据解析为请求对象

### Nonce生成规则

**要求**:
- 每次请求必须唯一
- 长度建议24-32个字符
- 使用密码学安全的随机数生成器

**示例**:
```python
import secrets
nonce = secrets.token_urlsafe(24)  # 生成32字符的URL安全随机字符串
```

### 时间戳验证规则

**服务端验证逻辑**:
```
当前时间 - replay_window <= 请求时间戳 <= 当前时间 + replay_window
```

**默认配置**: `replay_window = 300秒` (±5分钟)

**注意事项**:
- 确保客户端时间与服务器时间同步
- 时间差超过5分钟会导致请求被拒绝
- 错误提示: "timestamp expired" 或 "invalid timestamp"

---

## API调用流程

### 登录认证流程

```
客户端                                服务端
  │                                    │
  │  1. 准备登录数据                   │
  │  {                                │
  │    project_uuid,                  │
  │    card_key,                      │
  │    hwid (可选)                    │
  │  }                                │
  │                                    │
  │  2. AES加密数据                    │
  │                                    │
  │  3. POST /api/auth/login          │
  │  {                                │
  │    timestamp,                     │
  │    nonce,                         │
  │    data (encrypted)               │
  │  } ──────────────────────────────>│
  │                                    │
  │                                    │  4. 验证时间戳
  │                                    │  5. 验证nonce
  │                                    │  6. 解密数据
  │                                    │  7. 验证卡密
  │                                    │  8. 检查HWID/IP
  │                                    │  9. 生成Token
  │                                    │
  │  10. 返回Token                     │
  │ <────────────────────────────────│
  │  {                                │
  │    code: 0,                       │
  │    data: {                        │
  │      token,                       │
  │      expire_at,                   │
  │      card: {...}                  │
  │    }                              │
  │  }                                │
  │                                    │
  │  11. 保存Token用于后续请求         │
  │                                    │
```

### Token使用说明

**获取Token**:
- 通过 `/api/auth/login` 登录成功后获得
- Token有效期由服务端配置决定（默认3600秒）

**使用Token**:
```http
Authorization: Bearer {token}
```

**适用接口**:
- `/api/heartbeat` - 心跳验证
- `/api/cloud-var/:key` - 获取云变量
- `/api/card/custom-data` - 更新专属信息
- `/api/project/info` - 获取项目信息

**示例**:
```python
headers = {"Authorization": f"Bearer {token}"}
response = requests.post(f"{server_url}/api/heartbeat", headers=headers)
```

### 心跳保活机制

**作用**:
- 验证Token有效性
- 保持会话活跃
- 检测账号状态变化

**推荐策略**:
```python
import threading
import time

def heartbeat_loop(client):
    """每30秒发送一次心跳"""
    while True:
        time.sleep(30)
        try:
            client.heartbeat()
        except Exception as e:
            print(f"心跳失败: {e}")
            # 重新登录或退出
            break

# 启动后台心跳线程
threading.Thread(target=heartbeat_loop, args=(client,), daemon=True).start()
```

**注意事项**:
- 心跳间隔应小于Token有效期
- 建议30-60秒发送一次
- 心跳失败应尝试重新登录

### 错误处理最佳实践

#### 通用错误码

| Code | 说明 | 处理方式 |
|------|------|----------|
| 0 | 成功 | 继续执行 |
| 400 | 请求参数错误 | 检查请求数据格式 |
| 401 | 未授权/认证失败 | 重新登录 |
| 404 | 资源不存在 | 检查请求路径/参数 |
| 500 | 服务器错误 | 稍后重试或联系管理员 |

#### 常见错误处理

```python
def handle_api_call(func):
    """API调用装饰器"""
    def wrapper(*args, **kwargs):
        try:
            response = func(*args, **kwargs)
            result = response.json()
            
            if result['code'] == 0:
                return result['data']
            elif result['code'] == 401:
                # Token过期，重新登录
                print("Token过期，重新登录...")
                # 执行重新登录逻辑
                raise Exception("需要重新登录")
            else:
                raise Exception(f"API错误: {result['message']}")
                
        except requests.RequestException as e:
            print(f"网络错误: {e}")
            raise
        except Exception as e:
            print(f"未知错误: {e}")
            raise
    
    return wrapper
```

---

## 多语言示例

### Python客户端

完整示例: [examples/python-client.py](examples/python-client.py)

**核心代码**:
```python
from Crypto.Cipher import AES
import base64
import requests

class NextKeyClient:
    def __init__(self, server_url, project_uuid, aes_key):
        self.server_url = server_url
        self.project_uuid = project_uuid
        self.aes_key = aes_key[:32].encode()  # 取前32字节
        self.token = None
    
    def login(self, card_key, hwid=""):
        # 详见示例文件
        pass
```

### Go客户端

完整示例: [examples/go-client.go](examples/go-client.go)

**核心代码**:
```go
type Client struct {
    ServerURL   string
    ProjectUUID string
    AESKey      []byte
    Token       string
}

func (c *Client) Login(cardKey, hwid string) error {
    // 详见示例文件
}
```

### 其他语言参考

**C#**:
```csharp
using System.Security.Cryptography;

// AES-GCM加密
using var aes = new AesGcm(keyBytes);
aes.Encrypt(nonce, plaintext, ciphertext, tag);
```

**JavaScript/Node.js**:
```javascript
const crypto = require('crypto');

function encrypt(key, plaintext) {
    const cipher = crypto.createCipheriv('aes-256-gcm', key, nonce);
    const encrypted = Buffer.concat([
        cipher.update(plaintext, 'utf8'),
        cipher.final()
    ]);
    const tag = cipher.getAuthTag();
    return Buffer.concat([nonce, tag, encrypted]).toString('base64');
}
```

---

## 常见问题FAQ

### 1. 加密失败排查

**问题**: 加密数据后服务端返回解密失败

**排查步骤**:
1. 检查AES密钥是否正确
   ```python
   print(f"密钥长度: {len(aes_key)}")  # 应该是32
   ```

2. 检查密钥格式转换
   ```python
   # 正确方式
   key_bytes = aes_key[:32].encode()
   
   # 错误方式
   # key_bytes = aes_key.encode()  # 长度会是64
   ```

3. 检查加密模式是否为GCM
   ```python
   cipher = AES.new(key, AES.MODE_GCM)  # 必须是GCM模式
   ```

4. 检查数据拼接顺序
   ```python
   # 正确顺序: nonce + tag + ciphertext
   encrypted = cipher.nonce + tag + ciphertext
   ```

### 2. 时间戳过期处理

**错误信息**: "timestamp expired" 或 "invalid timestamp"

**原因**: 客户端与服务端时间差超过5分钟

**解决方案**:
1. 同步系统时间
   ```bash
   # Windows
   w32tm /resync
   
   # Linux
   sudo ntpdate time.windows.com
   ```

2. 检查时区设置
   ```python
   import time
   print(f"当前时间戳: {int(time.time())}")
   ```

3. 调整服务端replay_window（临时方案）
   ```yaml
   security:
     replay_window: 600  # 增加到10分钟
   ```

### 3. Nonce重放错误

**错误信息**: "nonce already used"

**原因**: 同一个nonce在时间窗口内被使用多次

**解决方案**:
1. 确保每次请求生成新的nonce
   ```python
   import secrets
   nonce = secrets.token_urlsafe(24)  # 每次都生成新的
   ```

2. 不要重复发送相同的请求体
   ```python
   # 错误做法
   request_body = {...}
   requests.post(url, json=request_body)  # 第一次
   requests.post(url, json=request_body)  # 第二次会失败
   
   # 正确做法
   def make_request():
       request_body = create_encrypted_request()  # 每次重新生成
       return requests.post(url, json=request_body)
   ```

### 4. HWID/IP限制说明

**HWID（硬件ID）限制**:

**配置位置**: 项目设置中的 `enable_hwid` 和卡密的 `max_hwid`

**行为**:
- `max_hwid = -1`: 不限制设备数量
- `max_hwid = 0`: 禁止任何设备登录
- `max_hwid = 1`: 仅允许1个设备（单机版）
- `max_hwid = N`: 允许N个不同设备

**示例**:
```python
# 绑定设备码
import uuid
hwid = str(uuid.getnode())  # 获取MAC地址作为设备码
client.login(card_key, hwid=hwid)
```

**IP限制**:

**配置位置**: 项目设置中的 `enable_ip` 和卡密的 `max_ip`

**行为**:
- `max_ip = -1`: 不限制IP数量
- `max_ip = 1`: 仅允许单一IP登录
- `max_ip = N`: 允许N个不同IP

**注意**:
- 动态IP用户可能频繁触发限制
- 建议配合HWID使用，提高安全性

### 5. 卡密激活与续期

**首次激活**:
```python
# 未激活的卡密首次登录会自动激活
result = client.login("NEWCARD-123456")
# activated: true
# expire_at: 当前时间 + duration
```

**过期检测**:
```python
from datetime import datetime

expire_at = result['card']['expire_at']
expire_time = datetime.fromisoformat(expire_at.replace('Z', '+00:00'))

if datetime.now() > expire_time:
    print("卡密已过期")
```

**续期方式**:
- 管理员在后台手动延期
- 使用新的卡密重新登录

### 6. 云变量使用场景

**适用场景**:
- 公告信息（版本更新、维护通知）
- 开关控制（功能开关、灰度发布）
- 配置参数（远程配置、动态参数）
- 防封验证（服务器状态检查）

**示例**:
```python
# 获取公告
notice = client.get_cloud_var("notice")
print(f"公告: {notice}")

# 检查功能开关
feature_enabled = client.get_cloud_var("new_feature") == "true"
if feature_enabled:
    # 启用新功能
    pass
```

### 7. 专属信息存储

**用途**: 存储用户级别的自定义数据

**示例**:
```python
# 保存用户数据
custom_data = {
    "user_level": 5,
    "points": 1000,
    "last_login": "2024-01-01"
}
client.update_custom_data(custom_data)

# 登录时获取
result = client.login(card_key)
saved_data = json.loads(result['card']['custom_data'])
print(f"用户等级: {saved_data['user_level']}")
```

### 8. 性能优化建议

**连接复用**:
```python
# 使用Session复用HTTP连接
import requests
session = requests.Session()
session.post(url, ...)  # 复用连接
```

**批量操作**:
```python
# 一次获取多个云变量（需要服务端支持）
# 当前版本需要逐个获取
vars = {}
for key in ['var1', 'var2', 'var3']:
    vars[key] = client.get_cloud_var(key)
```

**缓存Token**:
```python
# 保存Token到文件，避免频繁登录
import json

def save_token(token, expire_at):
    with open('token.json', 'w') as f:
        json.dump({'token': token, 'expire_at': expire_at}, f)

def load_token():
    try:
        with open('token.json', 'r') as f:
            data = json.load(f)
            # 检查是否过期
            return data
    except:
        return None
```

---

## 测试工具

我们提供了图形化测试工具，帮助您快速验证对接：

```bash
cd tools
pip install -r requirements.txt
python gui-test-client.py
```

详见: [tools/README.md](../tools/README.md)

---

## 技术支持

- **文档**: [API文档](API.md) | [部署指南](DEPLOY.md)
- **示例代码**: [examples/](examples/)
- **问题反馈**: [GitHub Issues](https://github.com/nextkey/nextkey/issues)

---

**更新时间**: 2024-01-01  
**文档版本**: 1.0.0

