package distributor

import (
	"context"
	"diploma/fuzz/distributor/mock"
	distributorpb "diploma/gen/distributor"
	"go.uber.org/mock/gomock"
	"testing"
)

//go:generate mockgen -source distributor.go -package mock -destination ./mock/mock.go

type DistributorServiceServer interface {
	distributorpb.DistributorServiceServer
}

type DistributorService struct {
	distributorpb.UnimplementedDistributorServiceServer
	mock *mock.MockDistributorServiceServer
}

func NewDistributorService(f *testing.F) *DistributorService {
	return &DistributorService{
		mock: mock.NewMockDistributorServiceServer(gomock.NewController(f)),
	}
}

func (s *DistributorService) Mock() *mock.MockDistributorServiceServer {
	return s.mock
}

func (s *DistributorService) DownloadFile(ctx context.Context, req *distributorpb.DownloadFileRequest) (*distributorpb.DownloadFileResponse, error) {
	return s.mock.DownloadFile(ctx, req)
}
