package main

import (
	"fmt"
	"log"
	"net"

	distributerpb "diploma/gen/distributor"
	uploaderpb "diploma/gen/uploader"
	"diploma/service/internal/mainservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run a service. It stores files in memory and gives them to the agent.
func main() {
	port := "50001"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	service := mainservice.NewMainService()

	uploaderpb.RegisterUploaderServiceServer(grpcServer, service)
	distributerpb.RegisterDistributorServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	fmt.Printf("Server is running on port %s.", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
