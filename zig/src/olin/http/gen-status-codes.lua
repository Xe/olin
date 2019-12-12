local dkjson = require "dkjson"

local fin, err = io.open("./status-codes.json", "r")
if err then
  error(err)
end

local content = fin:read "*a"
local status_codes, _, err = dkjson.decode(content, 1, nil)
if err then
  error(err)
end

local fout, err = io.open("./status_codes.zig", "w")
if err then
  error(err)
end

fout:write([[
/// StatusCode is a HTTP status code.
pub const StatusCode = enum(u32) {
]])

for i, v in pairs(status_codes) do
  local found = string.find(v["code"], "xx")
  if not found then
    local phrase = v["phrase"]
    local name = string.gsub(phrase, " ", "")
    name = string.gsub(name, "-", "")
    name = string.gsub(name, "'", "")
    local desc = string.sub(v["description"], 2, -2)
    local decl = string.format([[
  /// %s %s
  %s = %s,

]], name, desc, name, v["code"])
    fout:write(decl)
  end
end

fout:write("};\n\n")
fout:close()

os.execute("zig fmt status_codes.zig")
