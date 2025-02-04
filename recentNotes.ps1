Get-ChildItem "$env:CLOUD_DIR/Catalog/*.md" | 
    Where-Object { $_.LastWriteTime -gt (Get-Date).AddDays(-7)} | 
    Sort-Object -Property LastWriteTime -Descending | 
    Select-Object -Property @{Name="LastWriteTime";Expression={$_.LastWriteTime.ToString("yyyy-MM-dd hh:mm tt")}}, Name |
    ConvertTo-Json -Depth 7
