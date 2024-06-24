package main

import (
	"enquiry-service/grpc-proto-files/users"
	"enquiry-service/handlers"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartGrpcUserServer(localApiConfig *handlers.LocalApiConfig) {
	lis, err := net.Listen("tcp", ":50004")
	if err != nil {
		log.Fatalf("failed to listen on gRPC server: %v", err)
		return
	}

	server := grpc.NewServer()

	users.RegisterUserServiceServer(server, &handlers.UserServer{
		LocalApiConfig: localApiConfig,
	})
	log.Printf("grpc[user-service] server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
