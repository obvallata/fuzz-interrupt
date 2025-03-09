package main

import (
	"diploma/keypoint/client"
	"diploma/keypoint/injection"
	"log"
	"time"
)

func main() {
	c := client.NewKeyPointClient(client.Config{URL: "http://127.0.0.1:1234"})

	if err := c.Enable("now", injection.Config{
		Type:  injection.TypeSleep,
		Sleep: &injection.SleepInjectionConfig{Duration: 1 * time.Second},
	}); err != nil {
		log.Fatal(err)
	}
}
