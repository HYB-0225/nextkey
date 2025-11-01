-- NextKey SDK 使用示例

-- 导入 NextKey SDK
local NextKey = require("init")

-- 创建客户端实例
local client = NextKey.createClient({
    base_url = "http://localhost:8080",  -- 你的服务器地址
    project_uuid = "Your-uuid",  -- 项目UUID
    encrypt_key = "Your-encryption-key"  -- 32字节加密密钥
})

-- 执行登录
local success, result = client:login(
    "TEST-123456",           -- 卡密
    "DEVICE-HWID-001",      -- 设备码（可选）
    nil                      -- IP（可选，不传则使用请求IP）
)

if success then
    print("\n✓ 登录成功!")
    print("Token: " .. result.token)
    print("过期时间: " .. result.expire_at)
    print("卡密类型: " .. result.card.card_type)
    print("是否激活: " .. tostring(result.card.activated))
    print("是否冻结: " .. tostring(result.card.frozen))
    print("时长(秒): " .. result.card.duration)
    
    -- 定期发送心跳（保持会话活跃）
    local hbSuccess, hbResult = client:heartbeat()
    if hbSuccess then
        print("\n✓ 心跳成功")
    else
        print("\n✗ 心跳失败: " .. hbResult)
    end
    
    -- 获取云变量
    local varSuccess, varResult = client:getCloudVar("notice")
    if varSuccess then
        print("\n✓ 云变量 'notice' = " .. varResult.value)
    else
        print("\n✗ 获取云变量失败: " .. varResult)
    end
    
    -- 更新专属信息
    local customData = NextKey.json.encode({
        last_login = os.date("%Y-%m-%d %H:%M:%S"),
        device = "Android 12"
    })
    local updateSuccess, updateResult = client:updateCustomData(customData)
    if updateSuccess then
        print("\n✓ 专属信息更新成功")
    else
        print("\n✗ 更新失败: " .. updateResult)
    end
    
else
    print("\n✗ 登录失败: " .. result)
end

