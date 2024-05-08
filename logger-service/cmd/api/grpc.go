package main

import (
	"context"
	"log"
	"log-service/data"
	"log-service/logs/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Model data.Models
}

func (l LogServer) WriteLog(ctx context.Context, request *logs.LogRequest) (*logs.LogResponse, error) {
	input := request.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	log.Println("log grpc :", logEntry)

	err := l.Model.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged!"}
	return res, err
}
