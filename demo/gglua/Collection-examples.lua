-- ==================== 工具函数 ====================
-- 将十六进制字符串解码为字节数据
local function hexDecode(hexStr)
    if not hexStr or #hexStr % 2 ~= 0 then
        return nil
    end
    
    local bytes = {}
    for i = 1, #hexStr, 2 do
        local byte = tonumber(hexStr:sub(i, i + 1), 16)
        if not byte then
            return nil
        end
        table.insert(bytes, string.char(byte))
    end
    
    return table.concat(bytes)
end

-- ==================== RC4加密类 ====================
local RC4 = {}

function RC4:new(key)
    local obj = {}
    setmetatable(obj, self)
    self.__index = self
    
    obj.S = {}
    -- 先尝试hex解码，失败则直接使用原字符串（匹配服务端逻辑）
    local decodedKey = hexDecode(key)
    obj.key = decodedKey or key
    obj:init()
    
    return obj
end

function RC4:init()
    for i = 0, 255 do
        self.S[i] = i
    end
    
    local j = 0
    local keylen = #self.key
    for i = 0, 255 do
        j = (j + self.S[i] + self.key:byte((i % keylen) + 1)) % 256
        self.S[i], self.S[j] = self.S[j], self.S[i]
    end
    
    self.i = 0
    self.j = 0
end

function RC4:crypt(data)
    local result = {}
    local S = {}
    
    for k = 0, 255 do
        S[k] = self.S[k]
    end
    
    local i, j = self.i, self.j
    
    for k = 1, #data do
        i = (i + 1) % 256
        j = (j + S[i]) % 256
        S[i], S[j] = S[j], S[i]
        
        local K = S[(S[i] + S[j]) % 256]
        result[k] = string.char(bit32.bxor(data:byte(k), K))
    end
    
    return table.concat(result)
end

-- ==================== Base64编码/解码 ====================
local base64 = {}
local b64chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'

function base64.encode(data)
    return ((data:gsub('.', function(x) 
        local r,b='',x:byte()
        for i=8,1,-1 do r=r..(b%2^i-b%2^(i-1)>0 and '1' or '0') end
        return r;
    end)..'0000'):gsub('%d%d%d?%d?%d?%d?', function(x)
        if (#x < 6) then return '' end
        local c=0
        for i=1,6 do c=c+(x:sub(i,i)=='1' and 2^(6-i) or 0) end
        return b64chars:sub(c+1,c+1)
    end)..({ '', '==', '=' })[#data%3+1])
end

function base64.decode(data)
    data = string.gsub(data, '[^'..b64chars..'=]', '')
    return (data:gsub('.', function(x)
        if (x == '=') then return '' end
        local r,f='',(b64chars:find(x)-1)
        for i=6,1,-1 do r=r..(f%2^i-f%2^(i-1)>0 and '1' or '0') end
        return r;
    end):gsub('%d%d%d?%d?%d?%d?%d?%d?', function(x)
        if (#x ~= 8) then return '' end
        local c=0
        for i=1,8 do c=c+(x:sub(i,i)=='1' and 2^(8-i) or 0) end
        return string.char(c)
    end))
end

-- ==================== JSON编码/解码 ====================
local json = {}

function json.encode(obj)
    local function escape_str(s)
        return s:gsub('\\', '\\\\'):gsub('"', '\\"'):gsub('\n', '\\n'):gsub('\r', '\\r'):gsub('\t', '\\t')
    end
    
    local function encode_value(val)
        local t = type(val)
        if t == "string" then
            return '"' .. escape_str(val) .. '"'
        elseif t == "number" then
            return tostring(val)
        elseif t == "boolean" then
            return tostring(val)
        elseif t == "table" then
            local is_array = true
            local count = 0
            for k, v in pairs(val) do
                count = count + 1
                if type(k) ~= "number" or k < 1 or k > count then
                    is_array = false
                    break
                end
            end
            
            if is_array and count == #val then
                local arr = {}
                for i = 1, #val do
                    table.insert(arr, encode_value(val[i]))
                end
                return "[" .. table.concat(arr, ",") .. "]"
            else
                local obj = {}
                for k, v in pairs(val) do
                    table.insert(obj, '"' .. tostring(k) .. '":' .. encode_value(v))
                end
                return "{" .. table.concat(obj, ",") .. "}"
            end
        elseif val == nil then
            return "null"
        end
    end
    return encode_value(obj)
end

function json.decode(str)
    if not str or str == "" then
        return nil
    end
    
    local pos = 1
    
    local function skip_whitespace()
        while pos <= #str and str:sub(pos, pos):match("%s") do
            pos = pos + 1
        end
    end
    
    local function decode_value()
        skip_whitespace()
        if pos > #str then return nil end
        
        local char = str:sub(pos, pos)
        
        if char == '"' then
            pos = pos + 1
            local result = ""
            while pos <= #str do
                char = str:sub(pos, pos)
                if char == '"' then
                    pos = pos + 1
                    return result
                elseif char == "\\" then
                    pos = pos + 1
                    char = str:sub(pos, pos)
                    if char == "n" then result = result .. "\n"
                    elseif char == "r" then result = result .. "\r"
                    elseif char == "t" then result = result .. "\t"
                    elseif char == "\\" then result = result .. "\\"
                    elseif char == '"' then result = result .. '"'
                    else result = result .. char end
                    pos = pos + 1
                else
                    result = result .. char
                    pos = pos + 1
                end
            end
        elseif char == "{" then
            pos = pos + 1
            local obj = {}
            skip_whitespace()
            if str:sub(pos, pos) == "}" then
                pos = pos + 1
                return obj
            end
            while pos <= #str do
                skip_whitespace()
                local key = decode_value()
                skip_whitespace()
                if str:sub(pos, pos) ~= ":" then break end
                pos = pos + 1
                local value = decode_value()
                obj[key] = value
                skip_whitespace()
                if str:sub(pos, pos) == "}" then
                    pos = pos + 1
                    return obj
                end
                if str:sub(pos, pos) == "," then
                    pos = pos + 1
                end
            end
        elseif char == "[" then
            pos = pos + 1
            local arr = {}
            skip_whitespace()
            if str:sub(pos, pos) == "]" then
                pos = pos + 1
                return arr
            end
            while pos <= #str do
                table.insert(arr, decode_value())
                skip_whitespace()
                if str:sub(pos, pos) == "]" then
                    pos = pos + 1
                    return arr
                end
                if str:sub(pos, pos) == "," then
                    pos = pos + 1
                end
            end
        elseif char == "t" then
            if str:sub(pos, pos + 3) == "true" then
                pos = pos + 4
                return true
            end
        elseif char == "f" then
            if str:sub(pos, pos + 4) == "false" then
                pos = pos + 5
                return false
            end
        elseif char == "n" then
            if str:sub(pos, pos + 3) == "null" then
                pos = pos + 4
                return nil
            end
        else
            local num_str = ""
            while pos <= #str and str:sub(pos, pos):match("[%d%.%-eE+]") do
                num_str = num_str .. str:sub(pos, pos)
                pos = pos + 1
            end
            if num_str ~= "" then
                return tonumber(num_str)
            end
        end
    end
    
    return decode_value()
end

-- 生成随机字符串
local function generateRandomString(length)
    local chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    local result = ''
    math.randomseed(os.time() + os.clock() * 1000000)
    for i = 1, length do
        local rand = math.random(1, #chars)
        result = result .. chars:sub(rand, rand)
    end
    return result
end

-- ==================== 卡密登录客户端 ====================
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
    print("========== 开始卡密登录 ==========")
    
    -- 生成nonce和timestamp
    local nonce = generateRandomString(32)
    local timestamp = os.time()
    print("生成Nonce: " .. nonce)
    
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
    
    print("请求数据: " .. json.encode(businessData))
    
    -- 构造内层数据结构（服务端期望的格式）
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
    
    print("请求URL: " .. url)
    
    local response = gg.makeRequest(url, headers, requestBody)
    
    if not response then
        print("请求失败: 无响应")
        return false, "网络请求失败"
    end
    
    print("响应状态码: " .. tostring(response.code))
    print("响应内容: " .. tostring(response.content))
    
    -- 检查HTTP状态码
    if response.code ~= 200 then
        return false, "HTTP错误: " .. tostring(response.code)
    end
    
    -- 解析响应
    local responseData = json.decode(response.content)
    
    if not responseData then
        return false, "响应解析失败"
    end
    
    -- 检查是否是错误响应（未加密格式）
    if responseData.code and responseData.code ~= 0 and not responseData.data then
        -- 这是一个业务错误响应，直接返回错误信息
        return false, responseData.message or "登录失败"
    end
    
    -- 验证nonce（只有成功响应才有nonce）
    if not responseData.nonce then
        return false, "响应格式错误：缺少nonce字段"
    end
    
    if responseData.nonce ~= nonce then
        print("Nonce不匹配! 发送: " .. nonce .. ", 接收: " .. tostring(responseData.nonce))
        return false, "Nonce验证失败，可能存在安全风险"
    end
    
    -- 解密响应数据
    if not responseData.data then
        return false, "响应格式错误：缺少data字段"
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
    
    print("========== 登录成功 ==========")
    print("Token: " .. self.token)
    print("卡密信息: " .. json.encode(self.card_data))
    
    return true, innerResponse.data
end

-- 心跳验证
function CardLogin:heartbeat()
    if not self.token then
        return false, "未登录，请先调用login()"
    end
    
    print("========== 发送心跳 ==========")
    
    local nonce = generateRandomString(32)
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
    
    print("心跳成功")
    return true, innerResponse.data
end

-- 获取云变量
function CardLogin:getCloudVar(key)
    if not self.token then
        return false, "未登录"
    end
    
    print("========== 获取云变量: " .. key .. " ==========")
    
    local nonce = generateRandomString(32)
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
    
    print("请求URL: " .. url)
    print("请求体: " .. json.encode(fullRequest))
    
    local response = gg.makeRequest(url, headers, json.encode(fullRequest))
    
    if not response then
        print("响应为空")
        return false, "请求失败: 无响应"
    end
    
    print("响应状态码: " .. tostring(response.code))
    print("响应内容: " .. tostring(response.content))
    
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
    
    print("云变量值: " .. innerResponse.data.value)
    return true, innerResponse.data
end

-- 更新卡密专属信息
function CardLogin:updateCustomData(custom_data)
    if not self.token then
        return false, "未登录"
    end
    
    print("========== 更新专属信息 ==========")
    
    local nonce = generateRandomString(32)
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
    
    print("更新成功")
    return true, innerResponse.data
end

-- 解绑HWID
function CardLogin:unbindHWID(hwid)
    if not self.token then
        return false, "未登录"
    end
    
    print("========== 解绑HWID: " .. hwid .. " ==========")
    
    local nonce = generateRandomString(32)
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
    
    print("解绑成功")
    return true, innerResponse.data
end

-- ==================== 使用示例 ====================
-- 创建客户端实例
local client = CardLogin:new({
    base_url = "http://localhost:8080",  -- 你的服务器地址
    project_uuid = "Your-uuid",  -- 项目UUID
    encrypt_key = "Your-encryption-key"  -- 32字节加密密钥（需要与服务器一致）
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
    local customData = json.encode({
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