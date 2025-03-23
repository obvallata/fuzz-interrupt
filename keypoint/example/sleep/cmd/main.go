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

	if err := c.EnableInjection("now", injection.Config{
		Type:  injection.TypeSleep,
		Sleep: &injection.SleepInjectionConfig{Duration: 1 * time.Second},
	}); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	if err := c.DisableInjection("now"); err != nil {
		log.Fatal(err)
	}
}
