package tailscale

import (
	"context"
	"os"
	"slices"

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

func Devices(ctx context.Context, filter *string) ([]tailscale.Device, error) {
	client, err := getClient()

	if err!=nil {
		return nil, err
	}

	 devices, err:=client.Devices(ctx)

	if err != nil {
		return nil, err
	}

	if filter == nil {

	return devices, nil
	}

	filtered := slices.DeleteFunc(devices, func(d tailscale.Device) bool {
		return !slices.Contains(d.Tags, *filter)
	})

	return filtered, nil
}

// func getServices(ctx context.Context) (error, error) {
// 	client, err := getClient()

// 	if err!=nil {
// 		return nil, err
// 	}

// 	return nil,nil
// }