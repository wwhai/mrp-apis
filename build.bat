set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build
scp .\mrpapiserver ubuntu@119.91.206.97:/home/ubuntu