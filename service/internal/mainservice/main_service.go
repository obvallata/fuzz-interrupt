package mainservice

import (
	"context"
	distributerpb "diploma/gen/distributor"
	uploaderpb "diploma/gen/uploader"
	"diploma/service/internal/distributor"
	"diploma/service/internal/uploader"
)

type MainService struct {
	uploaderpb.UnimplementedUploaderServiceServer
	distributerpb.UnimplementedDistributorServiceServer
	uploader    *uploader.UploaderService
	distributor *distributor.DistributorService
}

func NewMainService() *MainService {
	s := MainService{
		uploader:    &uploader.UploaderService{},
		distributor: &distributor.DistributorService{},
	}
	s.uploader.InitUploaderService()
	s.distributor.InitDistributorService(s.uploader)
	return &s
}

func (s *MainService) UploadFile(ctx context.Context, req *uploaderpb.UploadFileRequest) (*uploaderpb.UploadFileResponse, error) {
	return s.uploader.UploadFile(ctx, req)
}

func (s *MainService) ListFiles(ctx context.Context, req *uploaderpb.ListFilesRequest) (*uploaderpb.ListFilesResponse, error) {
	return s.uploader.ListFiles(ctx, req)
}

func (s *MainService) GetFileInfo(ctx context.Context, req *uploaderpb.GetFileInfoRequest) (*uploaderpb.GetFileInfoResponse, error) {
	return s.uploader.GetFileInfo(ctx, req)
}

func (s *MainService) SetFileActive(ctx context.Context, req *uploaderpb.SetFileActiveRequest) (*uploaderpb.SetFileActiveResponse, error) {
	return s.uploader.SetFileActive(ctx, req)
}

func (s *MainService) GetFileActive(ctx context.Context, req *uploaderpb.GetFileActiveRequest) (*uploaderpb.GetFileActiveResponse, error) {
	return s.uploader.GetFileActive(ctx, req)
}

func (s *MainService) DownloadFile(ctx context.Context, req *distributerpb.DownloadFileRequest) (*distributerpb.DownloadFileResponse, error) {
	return s.distributor.DownloadFile(ctx, req)
}
