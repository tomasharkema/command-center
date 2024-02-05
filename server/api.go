package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	servicelist "github.com/darylblake/go-systemd-servicelist"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func mainIpv4Address(device tailscale.Device) string {
	for _, address := range device.Addresses {
		if !strings.Contains(address, ":") {
			return address
		}
	}

	return ""
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := servicelist.CollectServiceInfo()
	if err != nil {
		logger.Infoln("error marshalling to json")
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Infoln("error marshalling to json")
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	response, err := Devices(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Infoln("error marshalling to json")
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
	w.WriteHeader(http.StatusOK)
}

func attachApi(r *mux.Router) {
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/services", servicesHandler)
	s.HandleFunc("/devices", devicesHandler)
	s.HandleFunc("/status", statusHandler)
}
