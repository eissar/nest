-- show a popup window
vim.cmd("PopupWindow")
vim.cmd([[
term iwr "http://localhost:1323/api/ping" -ConnectionTimeoutSeconds 1 | Select StatusCode, Content && Invoke-WebRequest "http://localhost:1323/api/server/close" -ConnectionTimeoutSeconds 10 | Select StatusCode, Content; Write-Host "Running server..."; wt.exe -d "$env:CLOUD_DIR\Code\go\web-dashboard" pwsh -c ./build.ps1
]])
vim.cmd("startinsert")
