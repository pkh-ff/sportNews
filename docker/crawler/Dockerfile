# build
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ARG GOPROXY=https://goproxy.io,https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go build -ldflags '-s -w' -o /app/app ./cmd/crawler/main.go

# run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /app .
COPY conf/app.conf.example.yaml ./app.conf.yaml

CMD ["./app"]