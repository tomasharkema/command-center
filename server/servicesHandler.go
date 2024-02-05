package server

import (
	"encoding/json"
	"net/http"

	servicelist "github.com/darylblake/go-systemd-servicelist"
	"github.com/google/logger"
)

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	data, err := servicelist.CollectServiceInfo()
	if err != nil {
		logger.Infoln("error marshalling to json")
		http.Error(w, err.Error(), 500)
		return
	}

	enc := json.NewEncoder(w)

	err = enc.Encode(data)

	if err != nil {
		logger.Infoln("error marshalling to json")
		http.Error(w, err.Error(), 500)
		return
	}

}
