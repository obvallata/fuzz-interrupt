package keypoint

import (
	"diploma/keypoint/storage/mem"
	"fmt"
	"os"

	"diploma/keypoint/storage"
)

var keyPointStorage storage.KeyPointStorage

func init() {
	// TODO: disable env flag

	keyPointStorage = mem.NewMemStorage()

	if s := os.Getenv("GOFAIL_HTTP"); len(s) > 0 {
		if err := serve(s); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
