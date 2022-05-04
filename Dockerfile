FROM golang:1.18.1-alpine3.15 AS builder

COPY . /github.com/kosdirus/tgcrypto/
WORKDIR /github.com/kosdirus/tgcrypto/
RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /github.com/kosdirus/tgcrypto/bin/bot .
EXPOSE 8081
CMD ["./bot"]