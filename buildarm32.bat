CGO_ENABLED=1 GOARM=7 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc go build -ldflags "-s -w -linkmode external -extldflags -static"