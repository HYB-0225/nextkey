-- 工具函数模块

local utils = {}

-- 将十六进制字符串解码为字节数据
function utils.hexDecode(hexStr)
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

-- 生成随机字符串
function utils.generateRandomString(length)
    local chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    local result = ''
    math.randomseed(os.time() + os.clock() * 1000000)
    for i = 1, length do
        local rand = math.random(1, #chars)
        result = result .. chars:sub(rand, rand)
    end
    return result
end

return utils

