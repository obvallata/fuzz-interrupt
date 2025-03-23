package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"diploma/fuzz/internal/distributor"
	"diploma/fuzz/internal/fuzz"
	distributerpb "diploma/gen/distributor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 1. Поднимаем фейковый сервис.
// 2. Провоцируем агента на pull с помощью сигналов и кладем в ответ нагенеренные данные.
func main() {
	service := distributor.NewDistributorService()
	go upServer(service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()
	fuzz.FuzzDownloadFile(ctx, service)
}

func upServer(service *distributor.DistributorService) {
	port := "50001"
	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	distributerpb.RegisterDistributorServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	fmt.Printf("Server is running on port %s.\n", port)
	err = grpcServer.Serve(lis)
	fmt.Println("Server is stopping")
	if err != nil {
		fmt.Println(err.Error())
	}
}
