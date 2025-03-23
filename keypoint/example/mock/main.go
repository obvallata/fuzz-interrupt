package main

import (
	"context"
	"diploma/keypoint"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// go build && GOFAIL_HTTP="127.0.0.1:1234" ./mock
func main() {
	sum := func(a, b int) (int, error) { return a + b, nil }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Finished.")
			return
		default:
			if result, err := keypoint.WithInject(context.Background(), "sum", sum)(10, 5); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Success: %d\n", result)
			}

			time.Sleep(3 * time.Second)
		}
	}
}
