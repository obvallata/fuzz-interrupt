package main

import (
	"context"
	"diploma/keypoint"
	"time"
)

// go build -a -x -gcflags="all=-N -l"&& ./breakpoint
// dlv attach --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <PID>
func main() {
	var (
		ctx = context.Background()
		sum = func(a, b int) int { return a + b }
	)

	for {
		keypoint.WithInject(ctx, "start", func() {})()

		result := keypoint.WithInject(ctx, "sum", sum)(10, 5)
		if result < 0 {
			keypoint.WithInject(ctx, "positive", func() {})()
		} else {
			keypoint.WithInject(ctx, "negative", func() {})()
		}

		keypoint.WithInject(ctx, "start", func() {})()

		time.Sleep(3 * time.Second)
	}
}
