package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"diploma/keypoint"
)

// go build && GOKEYPOINT_HTTP="127.0.0.1:1234" ./notifications
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	sum := func(a, b int) int { return a + b }

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Finished.")
			return
		default:
			keypoint.State(ctx, "start")

			a, b := rand.Int(), rand.Int()
			keypoint.State(ctx, "numbers were generated", keypoint.WithData(map[string]any{
				"a": a,
				"b": b,
			}))

			result := keypoint.WithInject(ctx, "sum", sum)(a, b)
			fmt.Printf("[%s] sum result: %d\n", time.Now().Format(time.TimeOnly), result)
			keypoint.State(ctx, "sum was calculated", keypoint.WithData(map[string]any{
				"sum": result,
			}))

			keypoint.State(ctx, "finish")

			fmt.Println()
			time.Sleep(3 * time.Second)
		}
	}
}
