-- NextKey SDK Lua Client 主入口文件
-- 统一导出所有模块

local NextKey = {}

-- 加载各个模块
NextKey.utils = require("utils")
NextKey.base64 = require("base64")
NextKey.json = require("json")
NextKey.RC4 = require("rc4")
NextKey.CardLogin = require("client")

-- 版本信息
NextKey.version = "1.0.0"

-- 创建客户端的便捷方法
function NextKey.createClient(config)
    return NextKey.CardLogin:new(config)
end

return NextKey

