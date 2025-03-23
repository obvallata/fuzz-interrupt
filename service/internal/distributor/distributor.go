package distributor

import (
	"context"
	distributorpb "diploma/gen/distributor"
	uploaderpb "diploma/gen/uploader"
	"diploma/service/internal/uploader"
	"fmt"
)

type DistributorService struct {
	fileStorage *uploader.UploaderService
}

func (s *DistributorService) InitDistributorService(fileStorage *uploader.UploaderService) {
	s.fileStorage = fileStorage
}
func (s *DistributorService) DownloadFile(ctx context.Context, req *distributorpb.DownloadFileRequest) (*distributorpb.DownloadFileResponse, error) {
	active, err := s.fileStorage.GetFileActive(ctx, &uploaderpb.GetFileActiveRequest{})
	if err != nil {
		return nil, fmt.Errorf("no file to download")
	}

	result, err := s.fileStorage.GetFileInfo(ctx, &uploaderpb.GetFileInfoRequest{Path: active.Path})
	if err != nil {
		return nil, fmt.Errorf("failed to get active file")
	}
	return &distributorpb.DownloadFileResponse{Path: result.Path, Content: result.Content}, nil
}
