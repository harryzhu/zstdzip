SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o dist/linux/zstdzip -ldflags "-w -s" main.go

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o dist/windows/zstdzip.exe -ldflags "-w -s" main.go

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o dist/darwin_amd64/zstdzip -ldflags "-w -s" main.go

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=arm64
go build -o dist/darwin_arm64/zstdzip -ldflags "-w -s" main.go

7z.exe a dist/zstdzip_windows.zip dist/windows/zstdzip.exe
7z.exe a dist/zstdzip_linux.zip dist/linux/zstdzip
7z.exe a dist/zstdzip_mac.zip dist/darwin_arm64/zstdzip
7z.exe a dist/zstdzip_mac_intel.zip dist/darwin_amd64/zstdzip
