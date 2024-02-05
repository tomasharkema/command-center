package server

import (
	"log"
	"net/http"
	"time"

	"github.com/akyoto/cache"
	"github.com/alecthomas/kingpin/v2"
	"github.com/google/logger"
	systemd "github.com/iguanesolutions/go-systemd"
)

var (
	DevicesCache = cache.New(time.Minute)
	ip           = kingpin.Flag("listen", "IP address to ping.").Short('l').Default(":3333").TCP()
)

func StartServer() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/services", servicesHandler)

	if err := systemd.NotifyReady(); err != nil {
		logger.Errorf("failed to notify ready to systemd: %v", err)
	}
	addr := (*ip).String()
	logger.Infoln("Listening at:", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
