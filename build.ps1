$sourcecode = ".\cmd\maestro\main.go"
$target = "build\maestro"
$version = "1.0.1"
# Windows, 64-bit
$env:GOOS = 'windows'; $env:GOARCH = 'amd64';             go build -o "$($target)-v$($version)-win-amd64.exe" -ldflags "-s -w" $sourcecode
# Linux, 64-bit
$env:GOOS = 'linux';   $env:GOARCH = 'amd64';             go build -o "$($target)-v$($version)-linux-amd64"   -ldflags "-s -w" $sourcecode
# Raspberry Pi
$env:GOOS = 'linux';   $env:GOARCH = 'arm'; $env:GOARM=6; go build -o "$($target)-v$($version)-linux-armv6"   -ldflags "-s -w" $sourcecode
# macOS
$env:GOOS = 'darwin';  $env:GOARCH = 'amd64';             go build -o "$($target)-v$($version)-macos-amd64"   -ldflags "-s -w" $sourcecode