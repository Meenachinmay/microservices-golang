package main

import (
	enquiries "enquiry-service/enquiries-grpc"
	"enquiry-service/handlers"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartGrpcServer(localApiConfig *handlers.LocalApiConfig) {
	lis, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatalf("failed to listen on gRPC server: %v", err)
		return
	}

	server := grpc.NewServer()

	enquiries.RegisterEnquiryServiceServer(server, &handlers.EnquiryServer{
		LocalApiConfig: localApiConfig,
	})
	log.Printf("grpc[enquiry-service] server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
