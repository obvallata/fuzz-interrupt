package main

import (
	"log"
	"time"

	"diploma/keypoint/client"
	"diploma/keypoint/injection"
	"diploma/keypoint/schema"
)

func main() {
	c := client.NewKeyPointClient(client.Config{URL: "http://127.0.0.1:1234"})

	if err := c.EnableMonitor(schema.EnableMonitorRequest{}); err != nil {
		log.Fatal(err)
	}
	defer c.DisableMonitor()

	if err := c.EnableInjection("sum", injection.Config{
		Type: injection.TypeMock,
		Mock: &injection.MockInjectionConfig{Outs: []any{
			69, nil,
		}},
	}); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
}
