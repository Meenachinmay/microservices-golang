package handlers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log-service/internal/database"
	"log-service/logs"
	"net"
	"time"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	LocalApiConfig *LocalApiConfig
}

// WriteLog Handler method to write log using grpc
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// insert data here
	newLog, err := l.LocalApiConfig.DB.InsertLog(ctx, database.InsertLogParams{
		ServiceName: input.ServiceName,
		LogData:     input.LogData,
	})
	if err != nil {
		res := &logs.LogResponse{Result: "failed to save log via gRPC:[WriteLogGRPCHandler]" + err.Error()}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "success, logged!" + newLog.LogData}
	return res, nil
}

func (l *LogServer) GetAllLogs(ctx context.Context, request *logs.GetAllLogsRequest) (*logs.GetAllLogsResponse, error) {
	logEntries, err := l.LocalApiConfig.DB.GetAllLogs(ctx)
	if err != nil {
		return nil, errors.New("failed to get logs from database:gRPC:[WriteLogGRPCHandler]" + err.Error())
	}

	var logsResponse []*logs.Log
	for _, logEntry := range logEntries {
		logsResponse = append(logsResponse, &logs.Log{
			Id:          logEntry.ID,
			ServiceName: logEntry.ServiceName,
			LogData:     logEntry.LogData,
			CreatedAt:   logEntry.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   logEntry.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &logs.GetAllLogsResponse{Logs: logsResponse}, nil
}

func GRPCListener(localApiConfig *LocalApiConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", "50001"))
	if err != nil {
		log.Fatalf("failed to listen on gRPC server: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{
		LocalApiConfig: localApiConfig,
	})

	log.Printf("gRPC server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve for gRPC server: %v", err)
	}
}
