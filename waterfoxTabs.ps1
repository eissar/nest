return @{"WaterfoxTabs" = (Get-Process -Name waterfox).count} | ConvertTo-Json -Depth 7 -AsArray
