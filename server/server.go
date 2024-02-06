package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/akyoto/cache"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	systemd "github.com/iguanesolutions/go-systemd"
)

var (
	DevicesCache = cache.New(time.Minute)
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infoln(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func StartServer(addr string, withJobs bool, ctx context.Context) {
	if withJobs {
		go StartJobs(ctx)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	attachApi(r)
	r.Use(loggingMiddleware)

	if err := systemd.NotifyReady(); err != nil {
		logger.Errorf("failed to notify ready to systemd: %v", err)
	}
	logger.Infoln("Listening at:", addr)

	log.Fatal(http.ListenAndServe(addr, r))
}
