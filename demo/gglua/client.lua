-- 卡密登录客户端模块

local RC4 = require("rc4")
local base64 = require("base64")
local json = require("json")
local utils = require("utils")

local CardLogin = {}

function CardLogin:new(config)
    local obj = {
        base_url = config.base_url or "http://localhost:8080",
        project_uuid = config.project_uuid,
        encrypt_key = config.encrypt_key or "default-key-32-bytes-length!!",
        token = nil,
        card_data = nil
    }
    setmetatable(obj, self)
    self.__index = self
    return obj
end

-- 加密数据
function CardLogin:encryptData(data)
    local jsonData = json.encode(data)
    local rc4 = RC4:new(self.encrypt_key)
    local encrypted = rc4:crypt(jsonData)
    return base64.encode(encrypted)
end

-- 解密数据
function CardLogin:decryptData(encryptedData)
    local decoded = base64.decode(encryptedData)
    local rc4 = RC4:new(self.encrypt_key)
    local decrypted = rc4:crypt(decoded)
    return json.decode(decrypted)
end

-- 卡密登录
function CardLogin:login(card_key, hwid, ip)
    -- print("========== 开始卡密登录 ==========")
    
    -- 生成nonce和timestamp
    local nonce = utils.generateRandomString(32)
    local timestamp = os.time()
    -- print("生成Nonce: " .. nonce)
    
    -- 准备业务数据
    local businessData = {
        project_uuid = self.project_uuid,
        card_key = card_key
    }
    
    if hwid then
        businessData.hwid = hwid
    end
    
    if ip then
        businessData.ip = ip
    end
    
    -- print("请求数据: " .. json.encode(businessData))
    
    -- 构造内层数据结构(服务端期望的格式)
    local innerData = {
        nonce = nonce,
        timestamp = timestamp,
        data = businessData
    }
    
    -- 加密内层数据
    local encryptedData = self:encryptData(innerData)
    
    -- 组装外层请求
    local fullRequest = {
        timestamp = timestamp,
        nonce = nonce,
        data = encryptedData
    }
    
    local requestBody = json.encode(fullRequest)
    
    -- 发送POST请求
    local url = self.base_url .. "/api/auth/login"
    local headers = {
        ['Content-Type'] = 'application/json'
    }
    
    -- print("请求URL: " .. url)
    
    local response = gg.makeRequest(url, headers, requestBody)
    
    if not response then
        print("请求失败: 无响应")
        return false, "网络请求失败"
    end
    
    -- print("响应状态码: " .. tostring(response.code))
    -- print("响应内容: " .. tostring(response.content))
    
    -- 检查HTTP状态码
    if response.code ~= 200 then
        return false, "HTTP错误: " .. tostring(response.code)
    end
    
    -- 解析响应
    local responseData = json.decode(response.content)
    
    if not responseData then
        return false, "响应解析失败"
    end
    
    -- 检查是否是错误响应(未加密格式)
    if responseData.code and responseData.code ~= 0 and not responseData.data then
        -- 这是一个业务错误响应,直接返回错误信息
        return false, responseData.message or "登录失败"
    end
    
    -- 验证nonce(只有成功响应才有nonce)
    if not responseData.nonce then
        return false, "响应格式错误:缺少nonce字段"
    end
    
    if responseData.nonce ~= nonce then
        print("Nonce不匹配! 发送: " .. nonce .. ", 接收: " .. tostring(responseData.nonce))
        return false, "Nonce验证失败,可能存在安全风险"
    end
    
    -- 解密响应数据
    if not responseData.data then
        return false, "响应格式错误:缺少data字段"
    end
    
    local decryptedData = self:decryptData(responseData.data)
    
    if not decryptedData then
        return false, "响应解密失败"
    end
    
    -- 解密后的数据结构: { timestamp: xxx, data: { code: 0, message: "success", data: {...} } }
    local innerResponse = decryptedData.data
    if not innerResponse then
        return false, "响应内层数据格式错误"
    end
    
    -- 检查业务状态码
    if innerResponse.code ~= 0 then
        return false, innerResponse.message or "登录失败"
    end
    
    -- 保存token和卡密信息
    self.token = innerResponse.data.token
    self.card_data = innerResponse.data.card
    
    -- print("========== 登录成功 ==========")
    -- print("Token: " .. self.token)
    -- print("卡密信息: " .. json.encode(self.card_data))
    
    return true, innerResponse.data
end

-- 心跳验证
function CardLogin:heartbeat()
    if not self.token then
        return false, "未登录,请先调用login()"
    end
    
    -- print("========== 发送心跳 ==========")
    
    local nonce = utils.generateRandomString(32)
    local timestamp = os.time()
    
    -- 构造内层数据
    local innerData = {
        nonce = nonce,
        timestamp = timestamp,
        data = {}
    }
    
    local encryptedData = self:encryptData(innerData)
    
    local fullRequest = {
        timestamp = timestamp,
        nonce = nonce,
        data = encryptedData
    }
    
    local url = self.base_url .. "/api/heartbeat"
    local headers = {
        ['Content-Type'] = 'application/json',
        ['Authorization'] = 'Bearer ' .. self.token
    }
    
    local response = gg.makeRequest(url, headers, json.encode(fullRequest))
    
    if not response or response.code ~= 200 then
        return false, "心跳请求失败"
    end
    
    local responseData = json.decode(response.content)
    
    if not responseData or responseData.nonce ~= nonce then
        return false, "心跳响应验证失败"
    end
    
    local decryptedData = self:decryptData(responseData.data)
    local innerResponse = decryptedData.data
    
    if not innerResponse or innerResponse.code ~= 0 then
        return false, innerResponse and innerResponse.message or "心跳失败"
    end
    
    -- print("心跳成功")
    return true, innerResponse.data
end

-- 获取云变量
function CardLogin:getCloudVar(key)
    if not self.token then
        return false, "未登录"
    end
    
    -- print("========== 获取云变量: " .. key .. " ==========")
    
    local nonce = utils.generateRandomString(32)
    local timestamp = os.time()
    
    local innerData = {
        nonce = nonce,
        timestamp = timestamp,
        data = {}
    }
    
    local encryptedData = self:encryptData(innerData)
    
    local fullRequest = {
        timestamp = timestamp,
        nonce = nonce,
        data = encryptedData
    }
    
    local url = self.base_url .. "/api/cloud-var/" .. key
    local headers = {
        ['Content-Type'] = 'application/json',
        ['Authorization'] = 'Bearer ' .. self.token
    }
    
    -- print("请求URL: " .. url)
    -- print("请求体: " .. json.encode(fullRequest))
    
    local response = gg.makeRequest(url, headers, json.encode(fullRequest))
    
    if not response then
        print("响应为空")
        return false, "请求失败: 无响应"
    end
    
    -- print("响应状态码: " .. tostring(response.code))
    -- print("响应内容: " .. tostring(response.content))
    
    if response.code ~= 200 then
        return false, "请求失败: HTTP " .. tostring(response.code)
    end
    
    local responseData = json.decode(response.content)
    
    if not responseData or responseData.nonce ~= nonce then
        print("Nonce不匹配或响应为空")
        return false, "响应验证失败"
    end
    
    local decryptedData = self:decryptData(responseData.data)
    local innerResponse = decryptedData.data
    
    if not innerResponse or innerResponse.code ~= 0 then
        return false, innerResponse and innerResponse.message or "获取失败"
    end
    
    -- print("云变量值: " .. innerResponse.data.value)
    return true, innerResponse.data
end

-- 更新卡密专属信息
function CardLogin:updateCustomData(custom_data)
    if not self.token then
        return false, "未登录"
    end
    
    -- print("========== 更新专属信息 ==========")
    
    local nonce = utils.generateRandomString(32)
    local timestamp = os.time()
    
    local innerData = {
        nonce = nonce,
        timestamp = timestamp,
        data = {
            custom_data = custom_data
        }
    }
    
    local encryptedData = self:encryptData(innerData)
    
    local fullRequest = {
        timestamp = timestamp,
        nonce = nonce,
        data = encryptedData
    }
    
    local url = self.base_url .. "/api/card/custom-data"
    local headers = {
        ['Content-Type'] = 'application/json',
        ['Authorization'] = 'Bearer ' .. self.token
    }
    
    local response = gg.makeRequest(url, headers, json.encode(fullRequest))
    
    if not response or response.code ~= 200 then
        return false, "请求失败"
    end
    
    local responseData = json.decode(response.content)
    
    if not responseData or responseData.nonce ~= nonce then
        return false, "响应验证失败"
    end
    
    local decryptedData = self:decryptData(responseData.data)
    local innerResponse = decryptedData.data
    
    if not innerResponse or innerResponse.code ~= 0 then
        return false, innerResponse and innerResponse.message or "更新失败"
    end
    
    -- print("更新成功")
    return true, innerResponse.data
end

-- 解绑HWID
function CardLogin:unbindHWID(hwid)
    if not self.token then
        return false, "未登录"
    end
    
    -- print("========== 解绑HWID: " .. hwid .. " ==========")
    
    local nonce = utils.generateRandomString(32)
    local timestamp = os.time()
    
    local innerData = {
        nonce = nonce,
        timestamp = timestamp,
        data = {
            project_uuid = self.project_uuid,
            card_key = self.card_data.card_key,
            hwid = hwid
        }
    }
    
    local encryptedData = self:encryptData(innerData)
    
    local fullRequest = {
        timestamp = timestamp,
        nonce = nonce,
        data = encryptedData
    }
    
    local url = self.base_url .. "/api/card/unbind"
    local headers = {
        ['Content-Type'] = 'application/json',
        ['Authorization'] = 'Bearer ' .. self.token
    }
    
    local response = gg.makeRequest(url, headers, json.encode(fullRequest))
    
    if not response or response.code ~= 200 then
        return false, "请求失败"
    end
    
    local responseData = json.decode(response.content)
    
    if not responseData or responseData.nonce ~= nonce then
        return false, "响应验证失败"
    end
    
    local decryptedData = self:decryptData(responseData.data)
    local innerResponse = decryptedData.data
    
    if not innerResponse or innerResponse.code ~= 0 then
        return false, innerResponse and innerResponse.message or "解绑失败"
    end
    
    -- print("解绑成功")
    return true, innerResponse.data
end

return CardLogin
