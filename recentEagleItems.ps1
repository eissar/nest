#requires -modules EagleCoolUtils
$WarningPreference = 'SilentlyContinue'
return Get-EagleItemsV2 -OrderBy CREATEDATE -Limit 10 | Select-Object lastModified, name, ext, id, isDeleted, tags, url, annotation | ConvertTo-Json -Depth 7 -AsArray
