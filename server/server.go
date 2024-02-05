package server

import (
	"log"
	"net/http"
	"time"

	"github.com/akyoto/cache"
	"github.com/google/logger"
	systemd "github.com/iguanesolutions/go-systemd"
)

var DevicesCache = cache.New(time.Minute)

func StartServer() {
	http.HandleFunc("/", homeHandler)
	if err := systemd.NotifyReady(); err != nil {
		logger.Errorf("failed to notify ready to systemd: %v", err)
	}
	log.Fatal(http.ListenAndServe(":3333", nil))
}
