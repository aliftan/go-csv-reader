@echo off
:: build.bat
echo Building CSV Reader Application...

:: Set environment variables
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

:: Build the application
go build -ldflags="-H windowsgui" -o csv-reader.exe

echo Build complete! csv-reader.exe has been created.
pause