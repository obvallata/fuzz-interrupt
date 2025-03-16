package main

import (
	"context"
	"diploma/keypoint"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// go build -a -x -gcflags="all=-N -l"&& ./breakpoint
// dlv attach --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <PID>
func main() {
	fmt.Printf("PID: %d\n", os.Getpid())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	var (
		sum     = func(a, b int) int { return a + b }
		product = func(a, b int) int { return a * b }
	)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Finished.")
			return
		default:
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
}
