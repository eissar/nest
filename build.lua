function string.join(separator, list)
	local result = ""
	for i, v in ipairs(list) do
		result = result .. tostring(v) -- Convert each element to a string
		if i < #list then
			result = result .. separator
		end
	end
	return result
end
local space = string.char(32) --space

-- show command results in a popup buffer.
vim.cmd("PopupWindow")

local stopCmd =
	[[ $ErrorActionPreference = "SilentlyContinue"; iwr "http://localhost:1323/api/ping" -ConnectionTimeoutSeconds 1 | Select StatusCode, Content && Invoke-WebRequest "http://localhost:1323/api/server/close" -ConnectionTimeoutSeconds 10 | Select StatusCode, Content; Write-Host "Running server...";]]

-- powershell
local runCmd = [[wt.exe -d "]] .. vim.uv.cwd() .. [[" pwsh -c ./build.ps1]]

-- NOTE: vim.cmd handles newlines as the start of a new ex command.
vim.cmd(string.join(space, {
	[[term]],
	stopCmd,
	-- run ./build.ps1 in new console tab.
	runCmd,
}))
