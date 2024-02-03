package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tailscale/tailscale-client-go/tailscale"
)

func main() {
	apiKey := os.Getenv("TAILSCALE_API_KEY")
	tailnet := os.Getenv("TAILSCALE_TAILNET")

log.Println(apiKey, tailnet)

	client, err := tailscale.NewClient(apiKey, tailnet)
	if err != nil {
		log.Fatalln(err)
	}

	// List all your devices
	devices, err := client.Devices(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	sinceDate := time.Now()

	for i, device := range devices {
		timeAgo := device.LastSeen.Sub(sinceDate)
		log.Println(i, device.ID, device.Hostname, device.Addresses, device.LastSeen, timeAgo)
	}
}

