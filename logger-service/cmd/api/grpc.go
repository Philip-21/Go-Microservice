package main

import (
	"context"
	"fmt"
	"log"
	"log-service/database"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

//this file receives requests from the logs grpc servers

type LogServer struct {
	//it ensures backwards compatitbility
	//UnimplementedLogServiceServer must be embedded
	//to have forward compatible implementations.
	logs.UnimplementedLogServiceServer
	Models database.Models
}

// this func writes a response back to the broker after its saved in the db
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
	//write a log
	//create a log entry
	Logentry := database.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	err := l.Models.LogEntry.Insert(Logentry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}
	//return Response
	res := &logs.LogResponse{Result: "Logged"}
	return res, nil
}

// a Grpc listener
func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen to gRPC:%v", err)
	}
	s := grpc.NewServer()
	/*RegisterService registers a service and its implementation to
	the concrete type implementing this interface.*/
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("grpc Server Started on Port %s", gRpcPort)

	/*Serve accepts incoming connections on the listener lis, creating a new ServerTransport and service goroutine for each.
	The service goroutines read gRPC requests and then call the registered handlers to reply to them.
	Serve returns when lis.Accept fails with fatal errors. lis will be closed when this method returns*/
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen to gRPC:%v", err)
	}
}
