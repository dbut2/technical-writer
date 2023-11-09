FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go build -o /bin/technical-writer

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/technical-writer ./technical-writer

ENTRYPOINT ["./technical-writer"]
