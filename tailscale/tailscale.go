package tailscale

import (
	"context"
	"os"

	"github.com/tailscale/tailscale-client-go/tailscale"
)

func getClient() (*tailscale.Client, error)  {

	apiKey := os.Getenv("TAILSCALE_API_KEY")
	tailnet := os.Getenv("TAILSCALE_TAILNET")

	client, err := tailscale.NewClient(apiKey, tailnet)

	if err!=nil {
		return nil, err
	}

	return client, nil
}

func Devices(ctx context.Context) ([]tailscale.Device, error) {
	client, err := getClient()

	if err!=nil {
		return nil, err
	}

	 devices, err:=client.Devices(ctx)

	if err != nil {
		return nil, err
	}

	return devices, nil
}

// func getServices(ctx context.Context) (error, error) {
// 	client, err := getClient()

// 	if err!=nil {
// 		return nil, err
// 	}

// 	return nil,nil
// }