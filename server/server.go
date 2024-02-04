package server

import (
	"log"
	"net/http"
	"time"

	"github.com/akyoto/cache"
)


var DevicesCache = cache.New(time.Minute)

func StartServer() {
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":3333", nil))
}
