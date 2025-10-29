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

客户端API（`/api/*`）需要使用AES-256-GCM加密通信。

### 请求格式

```json
{
  "timestamp": 1698505200,
  "nonce": "随机32字符串",
  "data": "Base64编码的AES加密数据"
}
```

### 加密流程

1. 准备请求数据（JSON格式）
2. 使用AES-256-GCM加密
3. Base64编码
4. 组装完整请求

### 解密流程

1. 服务端验证timestamp（±5分钟）
2. 验证nonce（防重放）
3. Base64解码
4. AES-256-GCM解密
5. 解析JSON数据

### 响应格式

客户端敏感接口（卡密登录、心跳、云变量、卡密专属信息）的响应会携带Nonce进行加密：

```json
{
  "nonce": "客户端请求时发送的nonce",
  "data": "Base64编码的AES加密响应数据"
}
```

**响应验证流程**:

1. 客户端发送请求时记录发送的nonce
2. 收到响应后，验证响应中的`nonce`字段是否与发送的一致
3. Base64解码`data`字段
4. 使用AES-256-GCM解密
5. 解析JSON获取实际的响应数据（包含`code`、`message`、`data`字段）

**安全性**: 此机制防止攻击者将旧的有效响应重放给新的请求，即使响应被抓包，也无法用于其他请求。

## 客户端 API

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

**注意**: 如果卡密已被冻结（`frozen: true`），登录将失败并返回错误信息。

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

**接口**: `GET /api/cloud-var/:key`

**需要认证**: 是

**需要加密**: 是（请求和响应）

**请求说明**: GET 请求也需要发送加密的请求体（包含 timestamp、nonce、data 字段）

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

**接口**: `GET /api/project/info`

**需要认证**: 是

**需要加密**: 是（请求和响应）

**请求说明**: GET 请求也需要发送加密的请求体（包含 timestamp、nonce、data 字段）

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
  "token": "JWT Token"
}
```

### 2. 项目管理

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
  "projects": [
    {
      "name": "项目1",
      "mode": "free"
    },
    {
      "name": "项目2",
      "mode": "card"
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

### 3. 卡密管理

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

### 4. 云变量管理

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

