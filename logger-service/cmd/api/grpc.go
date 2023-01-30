package main

import (
	"context"
	"log-service/database"
	"log-service/logs"
)

//this file receives requests from the logs grpc servers

type LogServer struct {
	//it ensures backwards compatitbility
	//UnimplementedLogServiceServer must be embedded
	//to have forward compatible implementations.
	logs.UnimplementedLogServiceServer
	Models database.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
}
