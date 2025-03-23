package server

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"

	"diploma/chaosmachine/internal/action"
	"diploma/keypoint/schema"
)

// Serve sets up server for notifications
func Serve(config Config, action action.Action) error {
	h := &handler{action: action}
	http.HandleFunc("POST /", h.acceptNotification)

	ln, err := net.Listen("tcp", config.URL)
	if err != nil {
		return err
	}

	go http.Serve(ln, nil)
	return nil
}

type Config struct {
	URL string `yaml:"url"`
}

type handler struct {
	action action.Action
}

func (h *handler) acceptNotification(w http.ResponseWriter, r *http.Request) {
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

	h.action.HandleNotification(notification)
	w.WriteHeader(http.StatusNoContent)
}
