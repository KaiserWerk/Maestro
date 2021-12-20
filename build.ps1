$sourcecode = ".\cmd\maestro\main.go"
$target = "build\maestro"
# Windows, 64-bit
$env:GOOS = 'windows'; $env:GOARCH = 'amd64';             go build -o "$($target)-win-amd64.exe" -ldflags "-s -w" $sourcecode
# Linux, 64-bit
$env:GOOS = 'linux';   $env:GOARCH = 'amd64';             go build -o "$($target)-linux-amd64"   -ldflags "-s -w" $sourcecode
# Raspberry Pi
$env:GOOS = 'linux';   $env:GOARCH = 'arm'; $env:GOARM=6; go build -o "$($target)-linux-armv6"   -ldflags "-s -w" $sourcecode
# macOS
$env:GOOS = 'darwin';  $env:GOARCH = 'amd64';             go build -o "$($target)-macos-amd64"   -ldflags "-s -w" $sourcecode