#!/bin/sh

mkdir -p Releases

# 【darwin/amd64】
echo "start build darwin/amd64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./Releases/m3u8-darwin-amd64 cmd/m3u8-downloader/main.go

# 【linux/amd64】
echo "start build linux/amd64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./Releases/m3u8-linux-amd64 cmd/m3u8-downloader/main.go

# 【windows/amd64】
echo "start build windows/amd64 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./Releases/m3u8-windows-amd64.exe cmd/m3u8-downloader/main.go

echo "Congratulations,all build success!!!"
