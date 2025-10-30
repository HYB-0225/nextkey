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
   - 在管理后台创建项目，获取 `project_uuid` 和 `encryption_key`
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
       aes_key="your-encryption-key-from-admin"
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

### 加密密钥获取

每个项目拥有独立的加密密钥。从管理后台获取项目的加密密钥：

1. 登录管理后台
2. 进入"项目管理"页面
3. 创建或编辑项目时，可以看到"加密密钥"字段
4. 点击复制按钮复制密钥

**示例密钥格式**:
```
632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037
```

### 密钥格式说明

**密钥格式**: 64字符十六进制字符串（HEX编码）

**客户端使用**:
- **方式1**: 直接使用前32字节（推荐）
  ```python
  encryption_key = "632005a33ebb7619c1efd3853c7109f1c075c7bb86164e35da72916f9d4ef037"
  key_bytes = encryption_key[:32].encode()  # 取前32个字符转字节
  ```

- **方式2**: HEX解码为32字节
  ```python
  import binascii
  encryption_key = "632005a33ebb7619c1efd3853c7109f1..."
  key_bytes = binascii.unhexlify(encryption_key)  # 32字节
  ```

**重要**: 每个项目使用独立的加密密钥，确保使用正确的密钥与对应项目通信。

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
encryption_key = os.environ.get('NEXTKEY_ENCRYPTION_KEY')
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
- `/api/card/unbind` - 解绑HWID

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

### HWID解绑流程

**功能说明**:
- 允许用户主动解绑已绑定的HWID
- 需要项目启用解绑功能（管理后台配置）
- 受冷却时间限制，防止频繁解绑
- 可配置解绑扣时作为惩罚机制

**解绑流程**:
```
客户端                                服务端
  │                                    │
  │  1. 准备解绑数据                     │
  │  {                                 │
  │    project_uuid,                   │
  │    card_key,                       │
  │    hwid (要解绑的设备码)             │
  │  }                                 │
  │                                    │
  │  2. POST /api/card/unbind          │
  │  (需要Token认证)                    │
  │ ──────────────────────────────>    │
  │                                    │
  │                                    │  3. 验证项目是否启用解绑
  │                                    │  4. 检查冷却时间
  │                                    │  5. 验证HWID是否已绑定
  │                                    │  6. 从列表中移除HWID
  │                                    │  7. 扣除时间（如配置）
  │                                    │  8. 记录解绑历史
  │                                    │
  │  9. 返回解绑结果                     │
  │ <────────────────────────────────  │
  │  {                                 │
  │    code: 0,                        │
  │    message: "解绑成功"              │
  │  }                                 │
```

**Python示例**:
```python
def unbind_hwid(client, card_key, hwid):
    """解绑HWID"""
    data = {
        "project_uuid": client.project_uuid,
        "card_key": card_key,
        "hwid": hwid
    }
    result, _, _ = client.make_encrypted_request("/api/card/unbind", data)
    return result

# 使用示例
try:
    result = unbind_hwid(client, "your-card-key", "device-hwid-001")
    if result['code'] == 0:
        print("解绑成功！现在可以在其他设备上使用此卡密")
    else:
        print(f"解绑失败: {result['message']}")
except Exception as e:
    print(f"解绑异常: {e}")
```

**重要提示**:
- 解绑需要先登录获取Token
- 默认冷却时间为24小时，请勿频繁解绑
- 如果项目配置了解绑扣时，会减少卡密剩余时间
- 解绑后该HWID可以重新绑定

**常见错误**:
- `"该项目未启用解绑功能"` - 联系管理员在项目配置中启用
- `"解绑冷却中，请等待 XXX 秒后再试"` - 需要等待冷却时间
- `"该设备未绑定到此卡密"` - 检查HWID是否正确
- `"卡密已冻结"` - 卡密被冻结，无法解绑

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

### C++ 客户端 SDK

NextKey 提供基于 Rust FFI 的现代 C++ 客户端封装。

**特性**:
- ✅ 现代 C++17 风格
- ✅ RAII 自动资源管理
- ✅ 异常处理机制
- ✅ 异步心跳支持
- ✅ 跨平台（Windows、Linux、Android）

**快速开始**:
```cpp
#include "NextKeyClient.hpp"
using namespace nextkey;

int main() {
    try {
        auto client = std::make_unique<NextKeyClient>(
            "http://localhost:8080",
            "project-uuid",
            "aes-key"
        );
        
        auto result = client->login("card-key", "device-001");
        std::cout << "Token: " << result.token << "\n";
        
        client->heartbeat();
        
        auto value = client->getCloudVar("notice");
        std::cout << "Notice: " << value << "\n";
        
        // 解绑HWID（可选）
        try {
            client->unbindHWID("card-key", "old-device-hwid");
            std::cout << "HWID解绑成功\n";
        } catch (const NextKeyException& e) {
            std::cerr << "解绑失败: " << e.what() << "\n";
        }
        
    } catch (const NextKeyException& e) {
        std::cerr << "Error: " << e.what() << "\n";
        return 1;
    }
    return 0;
}
```

**完整文档**: [demo/cpp-client/README.md](../demo/cpp-client/README.md)

**编译与使用**:
```bash
# 复制 Rust 静态库
cp ../rust-sdk/target/release/libnextkey_sdk.a demo/cpp-client/lib/

# CMake 编译
cd demo/cpp-client
mkdir build && cd build
cmake ..
cmake --build . --config Release
```

### Rust SDK

原生 Rust 实现，支持编译为静态库、动态库或独立可执行文件。

**特性**:
- ✅ 零成本抽象
- ✅ 内存安全
- ✅ 跨平台编译（包括 Android ARM64）
- ✅ C FFI 支持（可供 C/C++ 调用）

**使用方式**:
```rust
use nextkey_sdk::client::NextKeyClient;

fn main() -> anyhow::Result<()> {
    let mut client = NextKeyClient::new(
        "http://localhost:8080",
        "project-uuid",
        "aes-key"
    )?;
    
    let result = client.login("card-key", Some("hwid"), None)?;
    println!("Token: {}", result.token);
    
    client.heartbeat()?;
    
    let value = client.get_cloud_var("notice")?;
    println!("Notice: {}", value);
    
    // 解绑HWID
    let unbind_result = client.unbind_hwid("card-key", "old-device-hwid")?;
    if unbind_result.code == 0 {
        println!("HWID解绑成功");
    }
    
    Ok(())
}
```

**构建**:
```bash
cd demo/rust-sdk

# 构建为静态库（供 C++ 使用）
cargo build --release

# 构建为可执行文件
cargo build --release --bin nextkey-cli

# 交叉编译 Android ARM64
cargo build --release --target aarch64-linux-android
```

**Cargo.toml 集成**:
```toml
[dependencies]
nextkey-sdk = { path = "../demo/rust-sdk" }
```

### Python客户端

完整示例: [demo/tools/gui-test-client.py](../demo/tools/gui-test-client.py)

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
        # 详见完整示例文件
        pass
```

完整文档: [demo/tools/README.md](../demo/tools/README.md)

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
# frozen: false
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

**卡密冻结**:
- 管理员可以在后台冻结卡密（`frozen: true`）
- 冻结后的卡密无法登录，会返回 401 错误
- 已登录的会话不受影响（直到 Token 过期）
- 解冻后可正常使用

**状态检查**:
```python
result = client.login("card-key")
if result['card']['frozen']:
    print("警告：此卡密已被冻结")
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
cd demo/tools
pip install -r requirements.txt
python gui-test-client.py
```

详见: [demo/tools/README.md](../demo/tools/README.md)

---

## 技术支持

- **文档**: [API文档](API.md) | [部署指南](DEPLOY.md)
- **SDK**: [C++ SDK](../demo/cpp-client/README.md) | [Rust SDK](../demo/rust-sdk/) | [Python 工具](../demo/tools/README.md)
- **问题反馈**: [GitHub Issues](https://github.com/HYB-0225/nextkey/issues)

---

**更新时间**: 2025-10-30  
**文档版本**: 1.0.1

