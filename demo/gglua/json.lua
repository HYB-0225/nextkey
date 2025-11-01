-- JSON编码/解码模块

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

return json

