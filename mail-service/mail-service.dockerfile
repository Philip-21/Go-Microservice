FROM golang:1.18-alpine AS builder 


RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o mailApp ./cmd/api

RUN chmod +x /app/mailApp



FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailApp /app


CMD ["/app/mailApp"]