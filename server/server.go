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
	logger.Infoln(ip)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/services", servicesHandler)

	if err := systemd.NotifyReady(); err != nil {
		logger.Errorf("failed to notify ready to systemd: %v", err)
	}

	log.Fatal(http.ListenAndServe(":3333", nil))
}
