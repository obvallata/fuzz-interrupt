package uploader

import (
	"context"
	uploaderpb "diploma/gen/uploader"
	"fmt"
	"maps"
	"slices"
	"sync"
)

type UploaderService struct {
	mu         sync.Mutex
	files      map[string]string
	activeFile string
}

func (s *UploaderService) InitUploaderService() {
	s.files = make(map[string]string)
}

func (s *UploaderService) UploadFile(ctx context.Context, req *uploaderpb.UploadFileRequest) (*uploaderpb.UploadFileResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[req.Name] = req.Content
	return &uploaderpb.UploadFileResponse{Success: true}, nil
}

func (s *UploaderService) ListFiles(ctx context.Context, req *uploaderpb.ListFilesRequest) (*uploaderpb.ListFilesResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &uploaderpb.ListFilesResponse{Filenames: slices.Collect(maps.Keys(s.files))}, nil
}

func (s *UploaderService) GetFileInfo(ctx context.Context, req *uploaderpb.GetFileInfoRequest) (*uploaderpb.GetFileInfoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	content, exists := s.files[req.Path]
	if !exists {
		return nil, fmt.Errorf("file not found")
	}
	return &uploaderpb.GetFileInfoResponse{Path: req.Path, Content: content}, nil
}

func (s *UploaderService) SetFileActive(ctx context.Context, req *uploaderpb.SetFileActiveRequest) (*uploaderpb.SetFileActiveResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.files[req.Path]
	if !exists {
		return &uploaderpb.SetFileActiveResponse{Success: false}, fmt.Errorf("file not found")
	}
	s.activeFile = req.Path
	return &uploaderpb.SetFileActiveResponse{Success: true}, nil
}

func (s *UploaderService) GetFileActive(ctx context.Context, req *uploaderpb.GetFileActiveRequest) (*uploaderpb.GetFileActiveResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeFile == "" {
		return nil, fmt.Errorf("active file is not set")
	}
	return &uploaderpb.GetFileActiveResponse{Path: s.activeFile}, nil
}
