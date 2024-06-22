package main

import (
	"context"
	enquiries "enquiry-service/enquiries-grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

type EnquiryServer struct {
	enquiries.UnimplementedEnquiryServiceServer
}

func (e *EnquiryServer) HandleCustomerEnquiry(ctx context.Context, request *enquiries.CustomerEnquiryRequest) (*enquiries.CustomerEnquiryResponse, error) {
	input := request.GetEnquiry()

	log.Printf("Processing customer enquiry:[DEBUG:HandleCustomerEnquiry]")

	// database insertion operation here
	log.Printf("inserting enquiry into database.%+v\n", input)

	//
	res := &enquiries.CustomerEnquiryResponse{
		Success: true,
		Message: "Successfully processed customer enquiry via gRPC.",
	}
	return res, nil
}

func StartGrpcServer() {
	lis, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatalf("failed to listen on gRPC server: %v", err)
		return
	}

	server := grpc.NewServer()

	enquiries.RegisterEnquiryServiceServer(server, &EnquiryServer{})
	log.Printf("grpc[enquiry-service] server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
