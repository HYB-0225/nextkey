-- RC4加密模块

local utils = require("utils")

local RC4 = {}

function RC4:new(key)
    local obj = {}
    setmetatable(obj, self)
    self.__index = self
    
    obj.S = {}
    -- 先hex解码，失败则直接使用原字符串（匹配服务端逻辑）
    local decodedKey = utils.hexDecode(key)
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

return RC4

