
echo ** Building linux/amd64 **
set GOARCH=amd64
set GOOS=linux
go build -v -o build/amd64/analyse  analyse.go

echo ** Building windows/amd64 **
set GOARCH=amd64
set GOOS=windows
go build -v -o build/amd64/analyse.exe analyse.go
