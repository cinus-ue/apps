@echo off
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

go build -ldflags="-s -w -H windowsgui" -o spy.exe

echo Finished Building
pause