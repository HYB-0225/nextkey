# NextKey API 文档

## 通用响应格式

所有 API 返回统一的 JSON 格式:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- `code`: 状态码，0表示成功，非0表示错误
- `message`: 提示信息
- `data`: 响应数据

## 加密通信

客户端API（`/api/*`）需要使用加密通信。服务端支持多种加密方案，每个项目可独立配置加密方案和密钥。

### 获取加密方案

**接口**: `GET /api/crypto/schemes`

**无需认证**: 此接口公开访问

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "scheme": "aes-256-gcm",
      "name": "AES-256-GCM",
      "description": "AES-256-GCM 加密（推荐）",
      "security_level": "secure",
      "performance": "medium",
      "is_deprecated": false
    },
    {
      "scheme": "chacha20-poly1305",
      "name": "ChaCha20-Poly1305",
      "description": "ChaCha20-Poly1305 加密",
      "security_level": "secure",
      "performance": "fast",
      "is_deprecated": false
    },
    {
      "scheme": "rc4",
      "name": "RC4",
      "description": "RC4 加密（兼容性，仅测试）",
      "security_level": "insecure",
      "performance": "fast",
      "is_deprecated": true
    },
    {
      "scheme": "xor",
      "name": "XOR",
      "description": "XOR 加密（兼容性，仅测试）",
      "security_level": "insecure",
      "performance": "fast",
      "is_deprecated": true
    },
    {
      "scheme": "custom-base64",
      "name": "自定义Base64",
      "description": "随机字符表Base64编码，简单混淆",
      "security_level": "insecure",
      "performance": "fast",
      "is_deprecated": false
    }
  ]
}
```

**使用说明**:
- 客户端应在首次连接时调用此接口获取服务端支持的加密方案
- 根据项目配置使用对应的加密方案和密钥
- 默认加密方案为 `aes-256-gcm`

### 请求格式

```json
{
  "timestamp": 1698505200,
  "nonce": "随机32字符串",
  "data": "Base64编码的加密数据"
}
```

`data` 解密后的内部结构:

```json
{
  "nonce": "同外层nonce",
  "timestamp": 1698505200,
  "data": {
    "...": "业务参数"
  }
}
```

### 加密流程

1. 准备业务请求数据（JSON格式）
2. 包装内部结构（嵌入nonce和timestamp）
3. 使用项目配置的加密方案加密内部结构
4. Base64编码
5. 组装完整请求

### 解密流程

1. 服务端验证timestamp（±5分钟）
2. 验证nonce（防重放）
3. Base64解码
4. 根据项目配置的加密方案解密数据
5. 解析内部结构并校验nonce/timestamp一致性
6. 提取内部业务数据

**注意**: 不同加密方案的具体解密细节由服务端自动处理，客户端只需使用正确的加密方案和密钥即可

### 响应格式

客户端敏感接口（卡密登录、心跳、云变量、卡密专属信息）的响应会携带Nonce进行加密：

```json
{
  "nonce": "客户端请求时发送的nonce",
  "data": "Base64编码的加密响应数据"
}
```

`data` 解密后的内部结构:

```json
{
  "nonce": "客户端请求时发送的nonce",
  "timestamp": 1698505201,
  "data": {
    "code": 0,
    "message": "success",
    "data": {}
  }
}
```

**响应验证流程**:

1. 客户端发送请求时记录发送的nonce
2. 收到响应后，验证响应中的`nonce`字段是否与发送的一致
3. Base64解码`data`字段
4. 使用项目配置的加密方案解密数据
5. 解析内部结构并校验nonce/timestamp
6. 提取内部业务响应数据（`code`、`message`、`data`）

**安全性**: 此机制防止攻击者将旧的有效响应重放给新的请求，即使响应被抓包，也无法用于其他请求。

## 客户端 API

### 0. 获取加密方案列表

**接口**: `GET /api/crypto/schemes`

**需要加密**: 否

**需要认证**: 否

**响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "scheme": "aes-256-gcm",
      "name": "AES-256-GCM",
      "description": "AES-256-GCM 加密（推荐）",
      "security_level": "secure",
      "performance": "medium",
      "is_deprecated": false
    },
    {
      "scheme": "chacha20-poly1305",
      "name": "ChaCha20-Poly1305",
      "description": "ChaCha20-Poly1305 加密",
      "security_level": "secure",
      "performance": "fast",
      "is_deprecated": false
    },
    {
      "scheme": "rc4",
      "name": "RC4",
      "description": "RC4 加密（兼容性，仅测试）",
      "security_level": "insecure",
      "performance": "fast",
      "is_deprecated": true
    },
    {
      "scheme": "xor",
      "name": "XOR",
      "description": "XOR 加密（兼容性，仅测试）",
      "security_level": "insecure",
      "performance": "fast",
      "is_deprecated": true
    },
    {
      "scheme": "custom-base64",
      "name": "自定义Base64",
      "description": "随机字符表Base64编码，简单混淆",
      "security_level": "insecure",
      "performance": "fast",
      "is_deprecated": false
    }
  ]
}
```

**使用场景**:
- 客户端初始化时获取服务端支持的加密方案
- 用于验证项目配置的加密方案是否被服务端支持
- 可用于动态选择加密方案（如果客户端支持多种方案）

### 1. 卡密登录

**接口**: `POST /api/auth/login`

**需要加密**: 是（请求和响应）

**请求参数**:
```json
{
  "project_uuid": "项目UUID",
  "card_key": "卡密",
  "hwid": "设备码（可选）",
  "ip": "IP地址（可选，不传则使用请求IP）"
}
```

**响应格式**:
```json
{
  "nonce": "请求时的nonce",
  "data": "加密的响应数据"
}
```

**解密后的响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "认证Token",
    "expire_at": "2024-01-01T00:00:00Z",
    "card": {
      "id": 1,
      "card_key": "xxx",
      "activated": true,
      "frozen": false,
      "duration": 2592000,
      "custom_data": "专属信息"
    }
  }
}
```

**注意**:
- 如果卡密已被冻结（`frozen: true`），登录将失败并返回错误信息。
- 免费模式下 `card` 字段可能为 `null`。

### 2. 心跳验证

**接口**: `POST /api/heartbeat`

**需要认证**: 是（Header: `Authorization: Bearer {token}`）

**需要加密**: 是（请求和响应）

**响应格式**:
```json
{
  "nonce": "请求时的nonce",
  "data": "加密的响应数据"
}
```

**解密后的响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "心跳成功"
  }
}
```

### 3. 获取云变量

**接口**: `GET /api/cloud-var/:key` 或 `POST /api/cloud-var/:key`

**需要认证**: 是

**需要加密**: 是（请求和响应）

**请求说明**: 
- 支持 GET 和 POST 两种方法（为兼容性保留 GET，建议使用 POST）
- 请求需要发送加密的请求体（包含 timestamp、nonce、data 字段）

**响应格式**:
```json
{
  "nonce": "请求时的nonce",
  "data": "加密的响应数据"
}
```

**解密后的响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "project_id": 1,
    "key": "变量名",
    "value": "变量值"
  }
}
```

### 4. 更新卡密专属信息

**接口**: `POST /api/card/custom-data`

**需要认证**: 是

**需要加密**: 是（请求和响应）

**请求参数**:
```json
{
  "custom_data": "任意字符串（可以是 JSON 格式，也可以是纯文本）"
}
```

**响应格式**:
```json
{
  "nonce": "请求时的nonce",
  "data": "加密的响应数据"
}
```

**解密后的响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "更新成功"
  }
}
```

### 5. 获取项目信息

**接口**: `GET /api/project/info` 或 `POST /api/project/info`

**需要认证**: 是

**需要加密**: 是（请求和响应）

**请求说明**: 
- 支持 GET 和 POST 两种方法（为兼容性保留 GET，建议使用 POST）
- 请求需要发送加密的请求体（包含 timestamp、nonce、data 字段）

**响应格式**:
```json
{
  "nonce": "请求时的nonce",
  "data": "加密的响应数据"
}
```

**解密后的响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "uuid": "项目UUID",
    "name": "项目名称",
    "version": "1.0.0",
    "update_url": "更新地址"
  }
}
```

### 6. 解绑HWID

**接口**: `POST /api/card/unbind`

**需要认证**: 否

**需要加密**: 是（请求和响应）

**请求参数**:
```json
{
  "project_uuid": "项目UUID",
  "card_key": "卡密",
  "hwid": "要解绑的设备码"
}
```

**响应格式**:
```json
{
  "nonce": "请求时的nonce",
  "data": "加密的响应数据"
}
```

**解密后的响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "解绑成功"
  }
}
```

**注意事项**:
- 不要求登录Token
- 需要项目启用解绑功能（`enable_unbind: true`）
- 如果项目设置了`unbind_verify_hwid: true`，则只能解绑已绑定的HWID
- 如果关闭验证（`unbind_verify_hwid: false`），仍需传入HWID参数以从列表中移除
- 解绑操作受冷却时间限制（默认86400秒/24小时）
- 如果配置了解绑扣时（`unbind_deduct_time`），会从卡密剩余时间中扣除对应秒数
- 解绑冻结的卡密会返回错误

**可能的错误码**:
- `400`: 参数错误或项目未启用解绑功能
- `400`: 解绑冷却中（错误信息会包含剩余冷却时间）
- `400`: 该设备未绑定到此卡密（当启用HWID验证时）
- `401`: 加密校验失败或认证失败

## 管理后台 API

### 1. 管理员登录

**接口**: `POST /admin/login`

**请求参数**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应数据**:
```json
{
  "access_token": "JWT Access Token",
  "refresh_token": "Refresh Token",
  "expires_in": 900
}
```

### 2. 刷新管理员Token

**接口**: `POST /admin/refresh`

**请求参数**:
```json
{
  "refresh_token": "Refresh Token"
}
```

**响应数据**:
```json
{
  "access_token": "JWT Access Token",
  "refresh_token": "Refresh Token",
  "expires_in": 900
}
```

### 3. 管理员注销

**接口**: `POST /admin/logout`

**需要认证**: 是（Header: `Authorization: Bearer {access_token}`）

**响应数据**:
```json
{
  "message": "注销成功"
}
```

### 4. 项目管理

#### 获取项目列表

**接口**: `GET /admin/projects`

**需要认证**: 是

#### 创建项目

**接口**: `POST /admin/projects`

**请求参数**:
```json
{
  "name": "项目名称",
  "mode": "free",
  "enable_hwid": true,
  "enable_ip": true,
  "version": "1.0.0",
  "token_expire": 3600,
  "description": "描述"
}
```

#### 按 UUID 获取项目

**接口**: `GET /admin/projects/:uuid`

**响应数据**:
```json
{
  "id": 1,
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "name": "项目名称",
  "mode": "free",
  "enable_hwid": true,
  "enable_ip": true,
  "version": "1.0.0",
  "update_url": "",
  "token_expire": 3600,
  "description": "描述"
}
```

#### 更新项目

**接口**: `PUT /admin/projects/:id`

**请求参数**:
```json
{
  "name": "新项目名称",
  "version": "1.0.1",
  "description": "更新描述"
}
```

#### 删除项目

**接口**: `DELETE /admin/projects/:id`

#### 批量创建项目

**接口**: `POST /admin/projects/batch`

**请求参数**:
```json
{
  "data": [
    {
      "name": "项目1",
      "mode": "free"
    },
    {
      "name": "项目2",
      "mode": "paid"
    }
  ]
}
```

#### 批量删除项目

**接口**: `DELETE /admin/projects/batch`

**请求参数**:
```json
{
  "ids": [1, 2, 3]
}
```

#### 更新项目加密方案

**接口**: `POST /admin/projects/:id/encryption`

**需要认证**: 是

**请求参数**:
```json
{
  "encryption_scheme": "chacha20-poly1305"
}
```

**响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "encryption_scheme": "chacha20-poly1305",
    "encryption_key": "新生成的64字符十六进制密钥"
  }
}
```

**注意事项**:
- 更新加密方案会自动生成新的加密密钥
- 更新后需要通知所有客户端使用新的加密方案和密钥
- 建议在无活跃用户时进行更新操作

### 5. 卡密管理

#### 获取卡密列表

**接口**: `GET /admin/cards`

**查询参数**:
- `project_id`: 项目ID（必需）
- `page`: 页码
- `page_size`: 每页数量
- `keyword`: 卡密关键词（模糊搜索）
- `card_type`: 卡密类型
- `activated`: 激活状态（"true"/"false"）
- `frozen`: 冻结状态（"true"/"false"）
- `note`: 备注（模糊搜索）
- `custom_data`: 专属信息（模糊搜索）
- `hwid`: 设备码（模糊搜索）
- `ip`: IP地址（模糊搜索）
- `start_time`: 创建时间范围-开始（格式：2024-01-01 00:00:00）
- `end_time`: 创建时间范围-结束（格式：2024-01-01 23:59:59）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1
  }
}
```

#### 批量生成卡密

**接口**: `POST /admin/cards`

**请求参数**:
```json
{
  "project_id": 1,
  "card_key": "自定义卡密（可选）",
  "prefix": "前缀",
  "suffix": "后缀",
  "count": 10,
  "duration": 2592000,
  "card_type": "normal",
  "max_hwid": -1,
  "max_ip": -1,
  "note": "备注"
}
```

#### 获取单个卡密

**接口**: `GET /admin/cards/:id`

**响应数据**:
```json
{
  "id": 1,
  "card_key": "TEST-123456",
  "project_id": 1,
  "activated": true,
  "frozen": false,
  "duration": 2592000,
  "expire_at": "2024-01-01T00:00:00Z",
  "note": "测试卡密",
  "card_type": "normal",
  "custom_data": "{}",
  "hwid_list": ["device-001"],
  "ip_list": ["192.168.1.1"],
  "max_hwid": -1,
  "max_ip": -1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 更新卡密

**接口**: `PUT /admin/cards/:id`

**请求参数**:
```json
{
  "note": "更新后的备注",
  "duration": 3600,
  "card_type": "vip",
  "max_hwid": 1,
  "max_ip": 1,
  "custom_data": "自定义数据",
  "hwid_list": ["device-001", "device-002"],
  "ip_list": ["192.168.1.1"]
}
```

#### 删除卡密

**接口**: `DELETE /admin/cards/:id`

#### 冻结卡密

**接口**: `PUT /admin/cards/:id/freeze`

**说明**: 冻结后的卡密无法登录，已登录的会话不受影响

**响应数据**:
```json
{
  "message": "冻结成功"
}
```

#### 解冻卡密

**接口**: `PUT /admin/cards/:id/unfreeze`

**响应数据**:
```json
{
  "message": "解冻成功"
}
```

#### 批量更新卡密

**接口**: `PUT /admin/cards/batch`

**请求参数**:
```json
{
  "ids": [1, 2, 3],
  "data": {
    "duration": 3600,
    "card_type": "vip",
    "note": "批量更新"
  }
}
```

#### 批量删除卡密

**接口**: `DELETE /admin/cards/batch`

**请求参数**:
```json
{
  "ids": [1, 2, 3]
}
```

#### 批量冻结卡密

**接口**: `PUT /admin/cards/batch/freeze`

**请求参数**:
```json
{
  "ids": [1, 2, 3]
}
```

#### 批量解冻卡密

**接口**: `PUT /admin/cards/batch/unfreeze`

**请求参数**:
```json
{
  "ids": [1, 2, 3]
}
```

### 6. 云变量管理

#### 获取云变量列表

**接口**: `GET /admin/cloud-vars`

**查询参数**:
- `project_id`: 项目ID

#### 设置云变量

**接口**: `POST /admin/cloud-vars`

**请求参数**:
```json
{
  "project_id": 1,
  "key": "变量名",
  "value": "变量值"
}
```

#### 删除云变量

**接口**: `DELETE /admin/cloud-vars/:id`

#### 批量设置云变量

**接口**: `POST /admin/cloud-vars/batch`

**请求参数**:
```json
{
  "data": [
    {
      "project_id": 1,
      "key": "notice",
      "value": "系统维护通知"
    },
    {
      "project_id": 1,
      "key": "feature_enabled",
      "value": "true"
    }
  ]
}
```

#### 批量删除云变量

**接口**: `DELETE /admin/cloud-vars/batch`

**请求参数**:
```json
{
  "ids": [1, 2, 3]
}
```

## 错误码

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/认证失败 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

