package tailscalehelper

import (
	"context"
	"slices"

	"github.com/alecthomas/kingpin/v2"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

var (
	apiKey  = kingpin.Flag("api-key", "Api Key").Envar("TAILSCALE_API_KEY").Required().String()
	tailnet = kingpin.Flag("tailnet", "Tailnet").Envar("TAILSCALE_TAILNET").Required().String()
)

func getClient() (*tailscale.Client, error) {

	client, err := tailscale.NewClient(*apiKey, *tailnet)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func Devices(ctx context.Context, filter *string) ([]tailscale.Device, error) {
	client, err := getClient()

	if err != nil {
		return nil, err
	}

	devices, err := client.Devices(ctx)

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
