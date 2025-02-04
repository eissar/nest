#requires -modules EagleCoolUtils
$WarningPreference = 'SilentlyContinue'
return Get-EagleItemsV2 -OrderBy CREATEDATE -Limit 10 | ConvertTo-Json -Depth 7 -AsArray
