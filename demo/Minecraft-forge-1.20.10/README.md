因为我不熟悉Java，避免代码混乱，我采取了压缩包的形式

## 重要提示：响应双层 Nonce 验证

### 安全机制说明
NextKey 服务端已升级为双层 nonce 验证机制，以防止响应数据被篡改。

**验证流程：**
1. **外层验证**：验证 HTTP 响应中的 `nonce` 字段与请求发送的 nonce 是否一致
2. **内层验证**：解密响应数据后，再次验证解密内容中的 `nonce` 字段

**响应数据结构：**
```json
// 外层（明文）
{
  "nonce": "客户端发送的nonce",
  "data": "加密的内容"
}

// 解密后的内层
{
  "nonce": "必须与外层nonce一致",
  "timestamp": 1234567890,
  "data": {
    "code": 0,
    "message": "success",
    "data": { ... }
  }
}
```

### 客户端对接要求

**必须实现的验证逻辑：**

```java
// 1. 验证外层 nonce
if (!response.getNonce().equals(requestNonce)) {
    throw new SecurityException("外层Nonce验证失败，可能遭受重放攻击！");
}

// 2. 解密响应数据
String decrypted = decrypt(response.getData());
JsonObject internal = parseJson(decrypted);

// 3. 验证内层 nonce（双重验证）
if (!internal.get("nonce").getAsString().equals(requestNonce)) {
    throw new SecurityException("内层Nonce验证失败，响应数据可能被篡改！");
}

// 4. 验证时间戳
long serverTimestamp = internal.get("timestamp").getAsLong();
long timeDiff = Math.abs(System.currentTimeMillis() / 1000 - serverTimestamp);
if (timeDiff > 300) {
    throw new SecurityException("响应时间戳异常，可能遭受离线攻击！");
}

// 5. 提取业务数据
JsonObject businessData = internal.getAsJsonObject("data");
```

### 错误处理建议

当验证失败时，应采取以下措施：
- **外层 nonce 不匹配**：记录日志，拒绝处理响应，可能是网络劫持或重放攻击
- **内层 nonce 不匹配**：记录日志并告警，响应数据已被篡改，禁止使用
- **时间戳异常**：可能是离线攻击或时钟不同步，建议同步系统时间后重试

### 适用接口
所有使用加密通信的接口都需要实现双层 nonce 验证：
- `/api/auth/login` - 卡密登录
- `/api/heartbeat` - 心跳保活
- `/api/cloud-var/{key}` - 获取云变量
- `/api/card/custom-data` - 更新卡密自定义数据
- `/api/card/unbind` - 解绑设备码
- `/api/project/info` - 获取项目信息