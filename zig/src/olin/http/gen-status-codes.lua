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
/// StatusCode is a HTTP status code. This matches the list at https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
pub const StatusCode = enum(u32) {
]])

for i, v in pairs(status_codes) do
  local found = string.find(v["code"], "xx")
  v["found"] = found
  if not found then
    local phrase = v["phrase"]
    local name = string.gsub(phrase, " ", "")
    name = string.gsub(name, "-", "")
    name = string.gsub(name, "'", "")
    v["name"] = name
    local desc = string.sub(v["description"], 2, -2)
    local decl = string.format([[
  /// %s %s
  %s = %s,

]], name, desc, name, v["code"])
    fout:write(decl)
  end
end

fout:write("};\n\n")

fout:write([[pub fn reasonPhrase(sc: StatusCode) []const u8 {
  return switch (sc) {]])

for i, v in pairs(status_codes) do
  if not v["found"] then
    fout:write(string.format([[StatusCode.%s => "%s"[0..],]], v["name"], v["phrase"]))
  end
end

fout:write([[else => "Unknown"[0..],]] .. "\n")
fout:write("};\n")
fout:write("}\n")
fout:close()

os.execute("zig fmt status_codes.zig")
