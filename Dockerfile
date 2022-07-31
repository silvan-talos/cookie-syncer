# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN GOOS=linux go build -o syncer cmd/main.go

CMD ["./syncer"]