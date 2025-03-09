package distributor

import (
	"context"
	"diploma/fuzz/internal/distributor/mock"
	"diploma/fuzz/internal/testreporter"
	distributorpb "diploma/gen/distributor"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source distributor.go -package MockService -destination ./MockService/MockService.go

type DistributorServiceServer interface {
	distributorpb.DistributorServiceServer
}

type DistributorService struct {
	distributorpb.UnimplementedDistributorServiceServer
	MockService *mock.MockDistributorServiceServer
}

func NewDistributorService() *DistributorService {
	return &DistributorService{
		MockService: mock.NewMockDistributorServiceServer(gomock.NewController(&testreporter.TestReporter{})),
	}
}

func (s *DistributorService) Mock() *mock.MockDistributorServiceServer {
	return s.MockService
}

func (s *DistributorService) DownloadFile(ctx context.Context, req *distributorpb.DownloadFileRequest) (*distributorpb.DownloadFileResponse, error) {
	return s.MockService.DownloadFile(ctx, req)
}
