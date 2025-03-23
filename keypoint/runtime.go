package keypoint

import (
	"fmt"
	"os"
	"sync/atomic"

	"diploma/keypoint/interaction"
	"diploma/keypoint/storage/mem"

	"diploma/keypoint/storage"
)

var (
	enabled         atomic.Bool
	keyPointStorage storage.KeyPointStorage
	notifier        interaction.Notifier
)

func init() {
	// TODO: disable env flag

	keyPointStorage = mem.NewMemStorage()

	if s := os.Getenv("GOKEYPOINT_HTTP"); len(s) > 0 {
		if err := serve(s); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
