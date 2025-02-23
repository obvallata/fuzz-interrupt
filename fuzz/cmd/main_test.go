package cmd

import (
	"context"
	"diploma/fuzz/distributor"
	"diploma/fuzz/distributor/mock"
	"diploma/fuzz/internal/caller"
	distributerpb "diploma/gen/distributor"
	"fmt"
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
	"testing"
)

func Fuzz(f *testing.F) {
	mockService := upServer(f)
	c := createCaller()
	f.Add([]byte("asdsadsa"))
	counter := 0
	f.Fuzz(func(t *testing.T, data []byte) {
		fz := fuzz.NewConsumer(data)
		createdString, err := fz.GetString()
		if err != nil {
			return
		}
		// mvp meaningless, further will represent constraints
		if len(createdString) < 2 {
			return
		}
		if createdString[0] != 'a' {
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)

		mockService.EXPECT().DownloadFile(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req *distributerpb.DownloadFileRequest) (*distributerpb.DownloadFileResponse, error) {
				defer wg.Done()

				return &distributerpb.DownloadFileResponse{
					Path:    fmt.Sprintf("/Users/ddr/fuzz-interrupt/exp/test_%d.txt", counter),
					Content: createdString,
				}, nil
			})

		err = c.Call()
		if err != nil {
			return
		}
		wg.Wait()

		counter++
	})
}

func upServer(f *testing.F) *mock.MockDistributorServiceServer {
	port := "778"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	service := distributor.NewDistributorService(f)

	distributerpb.RegisterDistributorServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	fmt.Printf("Server is running on port %s.", port)
	go grpcServer.Serve(lis)

	return service.Mock()
}

func createCaller() caller.Caller {
	return &caller.SignalCaller{}
}
