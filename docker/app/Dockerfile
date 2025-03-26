# build
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


#http://host.docker.internal:9011/api/v1/resource
COPY . .

RUN go build -ldflags '-s -w' -o /app/app ./cmd/app/main.go

# run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /app .
#COPY conf/app.conf.example.yaml ./app.conf.yaml

CMD ["./app"]
