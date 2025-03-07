package keypoint

import (
	"diploma/keypoint/injection"
	"encoding/json"
	"net"
	"net/http"
)

import (
	"io/ioutil"
)

func serve(host string) error {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	go http.Serve(ln, &HttpHandler{})
	return nil
}

// HttpHandler is used to handle keypoint Update/Disable requests
type HttpHandler struct{}

func (*HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	if len(key) == 0 || key[0] != '/' {
		http.Error(w, "malformed request URI", http.StatusBadRequest)
		return
	}
	key = key[1:]

	switch {
	// update keypoint
	case r.Method == "PUT":
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

		if err := keyPointStorage.UpdateInjectionConfig(key, config); err != nil {
			http.Error(w, "failed to update keypoint "+string(key), http.StatusBadRequest)
			return
		}

	// disable keypoint
	case r.Method == "DELETE":
		if err := keyPointStorage.Disable(key); err != nil {
			http.Error(w, "failed to disable keypoint "+err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Add("Allow", "DELETE")
		w.Header().Set("Allow", "PUT")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
