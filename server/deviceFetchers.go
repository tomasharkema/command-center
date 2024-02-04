package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/logger"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func fetchDeviceInfo(name string, ctx context.Context) (result *DeviceInfo, err error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	url := fmt.Sprintf("https://%s/webhook/hooks/info-json", name)
	logger.Infof("Fetch url for: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Errorln("Fetch error:", err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorln("Fetch error:", err)
		return nil, err
	}

	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)

	var deviceInfoResponse DeviceInfoResponse
	err = dec.Decode(&deviceInfoResponse)
if err!=nil {
	logger.Errorln("Fetch error:", err)
	return nil, err
}

	info := &DeviceInfo{
		Time: time.Now(),
		Response: &deviceInfoResponse,
	}

	return info, nil
}

const fetchDevicesInfoKey = "fetchDevicesInfoKey" 
func fetchDevicesInfo(devices []tailscale.Device, ctx context.Context) []*DeviceInfo {

	res, found := DevicesCache.Get(fetchDevicesInfoKey)
	if  found {
		return res.([]*DeviceInfo)
	}

	results := make([]*DeviceInfo, len(devices))
	var wg sync.WaitGroup

	for index, value := range devices {
		wg.Add(1)
		go func(index int, device tailscale.Device) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()

			info, err := fetchDeviceInfo(device.Name, ctx)
			if err != nil {
				results[index] = nil
				return
			}

			results[index] = info
		}(index, value)
	}
	wg.Wait()
	DevicesCache.Set(fetchDevicesInfoKey, results, time.Minute)
	return results
}