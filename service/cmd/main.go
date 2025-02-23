package main

import (
	distributerpb "diploma/gen/distributor"
	uploaderpb "diploma/gen/uploader"
	"diploma/service/internal/mainservice"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	port := "778"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	service := mainservice.NewMainService()

	uploaderpb.RegisterUploaderServiceServer(grpcServer, service)
	distributerpb.RegisterDistributorServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	fmt.Printf("Server is running on port %s.", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
