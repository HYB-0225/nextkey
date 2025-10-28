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

## 客户端 API

### 1. 卡密登录

**接口**: `POST /api/auth/login`

**需要加密**: 是

**请求参数**:
```json
{
  "project_uuid": "项目UUID",
  "card_key": "卡密",
  "hwid": "设备码（可选）",
  "ip": "IP地址（可选，不传则使用请求IP）"
}
```

**响应数据**:
```json
{
  "token": "认证Token",
  "expire_at": "2024-01-01T00:00:00Z",
  "card": {
    "id": 1,
    "card_key": "xxx",
    "activated": true,
    "duration": 2592000,
    "custom_data": "专属信息"
  }
}
```

### 2. 心跳验证

**接口**: `POST /api/heartbeat`

**需要认证**: 是（Header: `Authorization: Bearer {token}`）

**响应数据**:
```json
{
  "message": "心跳成功"
}
```

### 3. 获取云变量

**接口**: `GET /api/cloud-var/:key`

**需要认证**: 是

**响应数据**:
```json
{
  "id": 1,
  "project_id": 1,
  "key": "变量名",
  "value": "变量值"
}
```

### 4. 更新卡密专属信息

**接口**: `POST /api/card/custom-data`

**需要认证**: 是

**请求参数**:
```json
{
  "custom_data": "JSON字符串"
}
```

### 5. 获取项目信息

**接口**: `GET /api/project/info`

**需要认证**: 是

**响应数据**:
```json
{
  "uuid": "项目UUID",
  "name": "项目名称",
  "version": "1.0.0",
  "update_url": "更新地址"
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

#### 更新项目

**接口**: `PUT /admin/projects/:id`

#### 删除项目

**接口**: `DELETE /admin/projects/:id`

### 3. 卡密管理

#### 获取卡密列表

**接口**: `GET /admin/cards`

**查询参数**:
- `project_id`: 项目ID
- `page`: 页码
- `page_size`: 每页数量

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

#### 删除卡密

**接口**: `DELETE /admin/cards/:id`

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

## 错误码

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/认证失败 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

