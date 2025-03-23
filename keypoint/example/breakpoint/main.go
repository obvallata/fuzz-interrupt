package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"diploma/keypoint"
)

// go build -a -x -gcflags="all=-N -l"&& GOKEYPOINT_HTTP="127.0.0.1:1234" ./breakpoint
// dlv attach --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <PID>
func main() {
	log.Printf("PID: %d\n", os.Getpid())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	for {
		fmt.Println()

		select {
		case <-ctx.Done():
			fmt.Println("Finished.")
			return
		case <-time.Tick(3 * time.Second):
			readFromFile(ctx)
		}
	}
}

func readFromFile(ctx context.Context) {
	f, err := keypoint.WithInject(ctx, "open", os.Open)("important_file.txt")
	if err != nil {
		log.Println(err)
		return
	}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(b))

	if err := f.Close(); err != nil {
		log.Println(err)
	}
}
