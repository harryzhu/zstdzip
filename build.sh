go build -o dist/macos_arm/zstdzip -ldflags "-w -s" main.go
zip dist/macos_arm/zstdzip_macos_arm.zip dist/macos_arm/zstdzip

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macos_intel/zstdzip -ldflags "-w -s" main.go
zip dist/macos_intel/zstdzip_macos_intel.zip dist/macos_intel/zstdzip


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux_amd64/zstdzip -ldflags "-w -s" main.go
zip dist/linux_amd64/zstdzip_linux_amd64.zip dist/linux_amd64/zstdzip


CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows_amd64/zstdzip.exe -ldflags "-w -s" main.go
zip dist/windows_amd64/zstdzip_windows_amd64.zip dist/windows_amd64/zstdzip.exe
