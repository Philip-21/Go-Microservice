#Building stage

#the base go image
FROM golang:1.18-alpine  AS builder

#run the command on the docker image we are building
RUN mkdir /app

#Copy files from the root of the directory and all location into the docker image (/app)

COPY . /app


#declaring the current working directoryinside the image
WORKDIR /app

#Build go code , brokerApp is the app name
RUN CGO_ENABLE=0 go build -o brokerApp ./cmd/api

#run the chmod command and add the executable flag
RUN chmod +x /app/brokerApp


#----------The Main running image--------

#Build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

#copying file from the builder stage, the files are copied to /app 
COPY --from=builder /app/brokerApp /app

#it build te applicatio into the broker app
CMD [ "/app/brokerApp" ] 

#run docker compose -d on terminal 