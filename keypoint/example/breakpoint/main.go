package main

import (
	"context"
	"diploma/keypoint"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// go build -a -x -gcflags="all=-N -l"&& ./breakpoint
// dlv attach --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <PID>
func main() {
	fmt.Printf("PID: %d\n", os.Getpid())

	var (
		ctx     = context.Background()
		sum     = func(a, b int) int { return a + b }
		product = func(a, b int) int { return a * b }
	)

	for {
		// Start flow with dummy func
		keypoint.WithInject(ctx, "start", func() {})()

		a, b := rand.Int(), rand.Int()

		result := keypoint.WithInject(ctx, "sum", sum)(a, b)
		fmt.Printf("[%s] sum result: %d\n", time.Now().Format(time.TimeOnly), result)

		result = keypoint.WithInject(ctx, "product", product)(a, b)
		fmt.Printf("[%s] product result: %d\n", time.Now().Format(time.TimeOnly), result)

		// Finish flow with dummy func
		keypoint.WithInject(ctx, "finish", func() {})()

		time.Sleep(3 * time.Second)
	}
}
