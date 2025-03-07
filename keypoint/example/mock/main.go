package main

import (
	"context"
	"diploma/keypoint"
	"fmt"
	"time"
)

// go build && GOFAIL_HTTP="127.0.0.1:1234" ./mock
func main() {
	sum := func(a, b int) (int, error) { return a + b, nil }

	for {
		if result, err := keypoint.WithInject(context.Background(), "sum", sum)(10, 5); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		} else {
			fmt.Printf("Success: %d\n", result)
		}

		time.Sleep(3 * time.Second)
	}
}
