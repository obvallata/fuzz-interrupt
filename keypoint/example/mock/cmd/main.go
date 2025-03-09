package main

import (
	"diploma/keypoint/client"
	"diploma/keypoint/injection"
	"log"
)

func main() {
	c := client.NewKeyPointClient(client.Config{URL: "http://127.0.0.1:1234"})

	if err := c.Enable("sum", injection.Config{
		Type: injection.TypeMock,
		Mock: &injection.MockInjectionConfig{Outs: []any{
			69, nil,
		}},
	}); err != nil {
		log.Fatal(err)
	}
}
