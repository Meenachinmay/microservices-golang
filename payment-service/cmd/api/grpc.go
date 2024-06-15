package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	payment "payment-service/grpc/proto"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
}

func (p *PaymentServer) ProcessPayment(ctx context.Context, request *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	input := request.GetNewPayment()

	// simulate payment processing logic here
	log.Printf("Processing payment for User %s and Card no %s\n", input.CardHolderName, input.CardNumber)

	// return a successful payment response
	res := &payment.PaymentResponse{
		Result: "payment successfully processed",
	}

	if input.CardNumber != "4242-4242-4242-4242" {
		res.Result = "Payment processing failed"
	}

	return res, nil
}

func startGRPCServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "50002"))
	if err != nil {
		log.Fatalf("failed to listen on gRPC server for payment-service: %v", err)
	}

	server := grpc.NewServer()
	payment.RegisterPaymentServiceServer(server, &PaymentServer{})

	log.Printf("starting gRPC server (payment-service) on port 50002\n")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC server for payment-service: %v", err)
	}
}
