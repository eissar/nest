$GOOSList = "windows","darwin" # "linux"
$GOARCHList = "amd64", "arm64", "386"

# windows
foreach ($arch in $GOARCHList) {
    $os = "windows"
    $env:GOOS = $os
    $env:GOARCH = $arch
    Write-Host "Building for $os/$arch..."

    $outputName = "build/nest-$os-$arch.exe"

    go build -ldflags -H=windowsgui -o $outputName 
    Write-Host "Built $outputName"
}

# darwin
foreach ($arch in $GOARCHList) {
    $os = "darwin"
    $env:GOOS = $os
    $env:GOARCH = $arch
    Write-Host "Building for $os/$arch..."

    $outputName = "build/nest-$os-$arch"

    go build -o $outputName
    Write-Host "Built $outputName"

    Write-Host "All builds completed."
}
