package main

import (
	"context"
	"diploma/keypoint"
	"fmt"
	"time"
)

func main() {
	sum := func(a, b int) int { return a + b }

	for {
		result := keypoint.WithInject(context.Background(), "sum", sum)(10, 5)
		fmt.Println(result, time.Now().Format(time.TimeOnly))

		time.Sleep(3 * time.Second)
	}
}
