#BUILDING STAGE
FROM golang:1.18-alpine  AS builder


RUN mkdir /app


COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o MailApp ./cmd/api

RUN chmod +x /app/MailApp

#Running

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/MailApp /app
COPY templates /templates


CMD ["/app/MailApp"]