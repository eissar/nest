
@("amd64", "arm64") | ForEach-Object { $$arch = $$_; \$env:GOOS='windows'; \$env:GOARCH=\$arch; go build -ldflags -H=windowsgui -o build/nest-windows-\$arch.exe; Write-Host 'Built build/nest-windows-\$arch.exe'; }
