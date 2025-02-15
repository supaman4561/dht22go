FROM golang:1.23.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum .
RUN go mod download


COPY *.go .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

CMD ["/app/main"]
