package main

func main() {
	// starting the gRPC server
	go startGRPCServer()

	select {}
}
