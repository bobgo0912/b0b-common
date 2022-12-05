build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app
docker-img:
	docker build . --build-arg B0B_ENV=dev -t bob/xx:1.0
docker-run:
	docker run -d --name settlement -e B0B_ENV=prod bob/xx:1.0