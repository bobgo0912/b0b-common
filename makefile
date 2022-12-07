BUILD_PATH=github.com/bobgo0912/b0b-common/internal/config/handle
VERSION=$(shell git describe --always --match "v[0-9]*" HEAD)
BUILD_INFO=-ldflags "-X $(BUILD_INFO_IMPORT_PATH).Cfg.Version=$(VERSION)"

build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app $(BUILD_INFO)
docker-img:
	docker build . --build-arg B0B_ENV=dev -t bob/xx:1.0
docker-run:
	docker run -d --name settlement -e B0B_ENV=prod bob/xx:1.0