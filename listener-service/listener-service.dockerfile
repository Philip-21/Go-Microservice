#BUILDING STAGE
FROM golang:1.18-alpine  AS builder


RUN mkdir /app


COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o listenerApp .

RUN chmod +x /app/listenerApp

#Running stage 

FROM alpine:latest

RUN mkdir /app


COPY --from=builder /app/listenerApp /app

CMD [ "/app/listenerApp"]