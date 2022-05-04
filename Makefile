.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t tgcrypto:v0.1 .

start-container:
	docker run --name tgcrypto -p 8081:8081 --env-file .env tgcrypto:v0.1