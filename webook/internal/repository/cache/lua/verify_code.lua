local key = KEYS[1]
local cntKey = key..":cnt"
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
-- 验证次数超额
if cnt == nil or cnt <= 0 then
    return -1
end
-- 字符串比较，没必要转数字了
local code = redis.call("get", key)
if code == expectedCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    return -2
end