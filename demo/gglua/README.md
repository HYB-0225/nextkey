# NextKey SDK - Lua 客户端

NextKey SDK 的模块化 Lua 实现，用于 GameGuardian (GG) 脚本环境。

## 模块结构

```
gglua/
├── init.lua           # 主入口文件
├── utils.lua          # 工具函数模块
├── base64.lua         # Base64 编解码模块
├── json.lua           # JSON 编解码模块
├── rc4.lua            # RC4 加密模块
├── client.lua         # 卡密登录客户端模块
└── README.md          # 本文件
```

## 快速开始

### 方式 1: 使用统一入口(推荐)

```lua
-- 导入 NextKey SDK
local NextKey = require("init")

-- 创建客户端实例
local client = NextKey.createClient({
    base_url = "http://your-server.com",
    project_uuid = "your-project-uuid",
    encrypt_key = "your-32-byte-encryption-key"
})

-- 执行登录
local success, result = client:login("CARD-KEY", "DEVICE-HWID")

if success then
    print("登录成功!")
    print("Token: " .. result.token)
    
    -- 发送心跳
    client:heartbeat()
    
    -- 获取云变量
    local ok, data = client:getCloudVar("notice")
    if ok then
        print("公告: " .. data.value)
    end
else
    print("登录失败: " .. result)
end
```

### 方式 2: 直接使用各个模块

```lua
-- 只导入需要的模块
local CardLogin = require("client")
local json = require("json")

-- 创建客户端
local client = CardLogin:new({
    base_url = "http://your-server.com",
    project_uuid = "your-project-uuid",
    encrypt_key = "your-encryption-key"
})

-- 使用客户端...
```

### 方式 3: 直接使用Collection-examples.lua(不推荐)

## API 文档

### CardLogin 类

#### 创建实例

```lua
local client = CardLogin:new({
    base_url = "http://localhost:8080",     -- 服务器地址
    project_uuid = "your-project-uuid",      -- 项目UUID
    encrypt_key = "your-encryption-key"      -- 加密密钥
})
```

#### login(card_key, hwid, ip)

卡密登录

- `card_key`: 卡密
- `hwid`: 设备码(可选)
- `ip`: IP地址(可选，不传则使用请求IP)

返回: `success, result`

```lua
local success, result = client:login("CARD-KEY", "DEVICE-HWID")
if success then
    -- result.token: 访问令牌
    -- result.card: 卡密信息
    -- result.expire_at: 过期时间
end
```

#### heartbeat()

发送心跳保持会话活跃

返回: `success, result`

```lua
local success, result = client:heartbeat()
```

#### getCloudVar(key)

获取云变量

- `key`: 变量名

返回: `success, result`

```lua
local success, result = client:getCloudVar("notice")
if success then
    print(result.value)  -- 变量值
end
```

#### updateCustomData(custom_data)

更新卡密专属信息

- `custom_data`: JSON字符串

返回: `success, result`

```lua
local customData = json.encode({
    last_login = os.date("%Y-%m-%d %H:%M:%S"),
    device = "Android 12"
})
local success, result = client:updateCustomData(customData)
```

#### unbindHWID(hwid)

解绑设备码

- `hwid`: 要解绑的设备码

返回: `success, result`

```lua
local success, result = client:unbindHWID("DEVICE-HWID")
```

## 工具模块

### utils

```lua
local utils = require("utils")

-- 十六进制解码
local bytes = utils.hexDecode("48656c6c6f")

-- 生成随机字符串
local randomStr = utils.generateRandomString(32)
```

### base64

```lua
local base64 = require("base64")

-- 编码
local encoded = base64.encode("Hello World")

-- 解码
local decoded = base64.decode(encoded)
```

### json

```lua
local json = require("json")

-- 编码
local jsonStr = json.encode({name = "test", value = 123})

-- 解码
local obj = json.decode(jsonStr)
```

### RC4

```lua
local RC4 = require("rc4")

-- 创建加密器
local rc4 = RC4:new("encryption-key")

-- 加密/解密
local encrypted = rc4:crypt("Hello World")
local decrypted = rc4:crypt(encrypted)  -- RC4 是对称加密
```

## 依赖关系

- `client.lua` 依赖: `rc4.lua`, `base64.lua`, `json.lua`, `utils.lua`
- `rc4.lua` 依赖: `utils.lua`
- 其他模块无外部依赖

## 注意事项

1. 确保 `encrypt_key` 与服务端配置一致
2. 定期调用 `heartbeat()` 保持会话活跃
3. 所有网络请求使用 `gg.makeRequest()` API
4. 需要 Lua 5.2+ 环境支持(bit32 库)

## 许可证

与 NextKey SDK 主项目保持一致

