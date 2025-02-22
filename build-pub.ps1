$GOOSList = "windows", "linux" # "darwin"
$GOARCHList = "amd64", "arm64", "386"

foreach ($os in $GOOSList) {
    foreach ($arch in $GOARCHList) {
        $outputName = "build/nest-$os-$arch"
        if ($os -eq "windows") {
            $outputName += ".exe"
        }
    
        Write-Host "Building for $os/$arch..."
        $env:GOOS = $os
        $env:GOARCH = $arch
        go build -o $outputName
        Write-Host "Built $outputName"
    }
}

Write-Host "All builds completed."
