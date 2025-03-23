package main

import (
	"context"
	"fmt"
	"time"

	"diploma/keypoint"
)

// go build && GOKEYPOINT_HTTP="127.0.0.1:1234" ./sleep
func main() {
	for {
		now := keypoint.WithInject(context.Background(), "now", time.Now)()
		fmt.Println(now)
	}
}
