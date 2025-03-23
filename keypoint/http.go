package keypoint

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"

	"diploma/keypoint/injection"
	"diploma/keypoint/interaction"
	"diploma/keypoint/schema"
	"diploma/keypoint/utils/ptr"
)

func serve(host string) error {
	http.HandleFunc("PUT /injection/{name}", enableInjection)
	http.HandleFunc("DELETE /injection/{name}", disableInjection)
	http.HandleFunc("POST /monitor/enable", enableMonitor)
	http.HandleFunc("POST /monitor/disable", disableMonitor)

	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	go http.Serve(ln, nil)

	return nil
}

func enableInjection(w http.ResponseWriter, r *http.Request) {
	injectionName := r.PathValue("name")

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed ReadAll in PUT", http.StatusBadRequest)
		return
	}

	// TODO: validate
	var config injection.Config
	if err := json.Unmarshal(v, &config); err != nil {
		http.Error(w, "failed Unmarshal in PUT", http.StatusBadRequest)
		return
	}

	if err := keyPointStorage.UpdateInjectionConfig(injectionName, config); err != nil {
		http.Error(w, "failed to update keypoint "+string(injectionName), http.StatusBadRequest)
	}
}

func disableInjection(w http.ResponseWriter, r *http.Request) {
	injectionName := r.PathValue("name")

	if err := keyPointStorage.Disable(injectionName); err != nil {
		http.Error(w, "failed to disable keypoint "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func enableMonitor(w http.ResponseWriter, r *http.Request) {
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed ReadAll in POST", http.StatusBadRequest)
		return
	}

	// TODO: validate
	var request schema.EnableMonitorRequest
	if err := json.Unmarshal(v, &request); err != nil {
		http.Error(w, "failed Unmarshal in POST", http.StatusBadRequest)
		return
	}

	if request.NotifierURL != nil {
		notifier = interaction.NewNotifierClient(interaction.NotifierConfig{URL: ptr.From(request.NotifierURL)})
	}
	enabled.Store(true)
}

func disableMonitor(w http.ResponseWriter, r *http.Request) {
	enabled.Store(false)
	notifier = nil
}
