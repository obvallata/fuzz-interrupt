package fuzz

import (
	"context"
	"diploma/fuzz/internal/caller"
	"diploma/fuzz/internal/distributor"
	distributerpb "diploma/gen/distributor"
	"encoding/base64"
	fuzz "github.com/google/gofuzz"
	"go.uber.org/mock/gomock"
	"sync"
)

func FuzzDownloadFile(ctx context.Context, service *distributor.DistributorService) {
	c := caller.NewSignalCaller()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fz := fuzz.New()

			var createdString string
			fz.Fuzz(&createdString)

			if !satisfyConstraints(createdString) {
				continue
			}

			var wg sync.WaitGroup
			wg.Add(1)
			service.Mock().EXPECT().DownloadFile(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, req *distributerpb.DownloadFileRequest) (*distributerpb.DownloadFileResponse, error) {
					defer wg.Done()

					return &distributerpb.DownloadFileResponse{
						Path:    "/Users/ddr/fuzz-interrupt/exp/" + createdString + ".txt",
						Content: createdString,
					}, nil
				})

			err := c.Call()
			if err != nil {
				return
			}
			wg.Wait()
		}
	}
}

func satisfyConstraints(generated string) bool {
	// mvp meaningless, further will represent real constraints
	if len(generated) < 2 {
		return false
	}
	if generated[0] != 'a' {
		return false
	}
	if _, err := base64.StdEncoding.DecodeString(generated); err != nil {
		return false
	}

	return true
}
