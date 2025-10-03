package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"diploma/keypoint/client"
	"diploma/keypoint/schema"
	"diploma/keypoint/utils/ptr"
	"github.com/kr/pretty"
)

const notifierURL = "127.0.0.1:1235"

var notifications = make(chan schema.NotifyRequest)

func main() {
	if err := serve(); err != nil {
		log.Fatal(err)
	}

	c := client.NewKeyPointClient(client.Config{URL: "http://127.0.0.1:1234"})

	if err := c.EnableMonitor(schema.EnableMonitorRequest{NotifierURL: ptr.T("http://" + notifierURL)}); err != nil {
		log.Fatal(err)
	}
	defer c.DisableMonitor()

	for i := 0; i < 20; i++ {
		notification := <-notifications
		pretty.Println(notification)
	}
}

func serve() error {
	http.HandleFunc("POST /", acceptNotification)

	ln, err := net.Listen("tcp", notifierURL)
	if err != nil {
		return err
	}

	go http.Serve(ln, nil)

	return nil
}

func acceptNotification(w http.ResponseWriter, r *http.Request) {
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed ReadAll in PUT", http.StatusBadRequest)
		return
	}

	// TODO: validate
	var notification schema.NotifyRequest
	if err := json.Unmarshal(v, &notification); err != nil {
		http.Error(w, "failed Unmarshal in PUT", http.StatusBadRequest)
		return
	}

	notifications <- notification

	w.WriteHeader(http.StatusNoContent)
}
