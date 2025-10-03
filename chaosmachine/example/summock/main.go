package main

import (
	"context"
	"fmt"
	"os"

	"diploma/keypoint"
)

// go build -a -x -gcflags="all=-N -l" && GOFAIL_HTTP="127.0.0.1:1234" ./summock
// dlv attach --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <PID>
func main() {
	fmt.Printf("PID: %d\n", os.Getpid())

	var (
		ctx = context.Background()
		sum = func(a, b int) int { return a + b }
	)

	for {
		keypoint.WithInject(ctx, "start", func() {})()

		result := keypoint.WithInject(ctx, "sum", sum)(10, 5)
		if result < 0 {
			keypoint.WithInject(ctx, "negative", func() {})()
		} else if result == 0 {
			keypoint.WithInject(ctx, "zero", func() {})()
		} else {
			keypoint.WithInject(ctx, "positive", func() {})()
		}

		keypoint.WithInject(ctx, "finish", func() {})()
	}
}
