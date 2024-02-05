package server

import (
	"log"
	"net/http"
	"time"

	"github.com/akyoto/cache"
	"github.com/alecthomas/kingpin/v2"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	systemd "github.com/iguanesolutions/go-systemd"
)

var (
	DevicesCache = cache.New(time.Minute)
	ip           = kingpin.Flag("listen", "IP address to ping.").Short('l').Default(":3333").TCP()
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logger.Infoln(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func StartServer() {

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	attachApi(r)
	r.Use(loggingMiddleware)

	if err := systemd.NotifyReady(); err != nil {
		logger.Errorf("failed to notify ready to systemd: %v", err)
	}
	addr := (*ip).String()
	logger.Infoln("Listening at:", addr)

	log.Fatal(http.ListenAndServe(addr, r))
}
