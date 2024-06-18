package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log-service/internal/database"
	"log-service/logs"
	"net"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	localApiConfig *LocalApiConfig
}

// WriteLog Handler method to write log using grpc
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// insert data here
	newLog, err := l.localApiConfig.DB.InsertLog(ctx, database.InsertLogParams{
		ServiceName: input.Name,
		LogData:     input.Data,
	})
	if err != nil {
		res := &logs.LogResponse{Result: "failed to save log via gRPC:[WriteLogGRPCHandler]" + err.Error()}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "success, logged!" + newLog.LogData}
	return res, nil
}

func (localApiConfig *LocalApiConfig) GRPCListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", "50001"))
	if err != nil {
		log.Fatalf("failed to listen on gRPC server: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{})

	log.Printf("gRPC server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve for gRPC server: %v", err)
	}
}
