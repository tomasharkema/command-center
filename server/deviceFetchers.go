package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/tomasharkema/go-nixos-menu/tailscalehelper"

	"github.com/google/logger"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"github.com/xeonx/timeago"
)

func fetchDeviceStatus(name string, ctx context.Context) (*string, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	url := fmt.Sprintf("http://%s:3333/api/services", name)
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
	var resultBuffer strings.Builder
	_, err = io.Copy(&resultBuffer, res.Body)
	if err != nil {
		logger.Errorln("Fetch error:", err)
		return nil, err
	}
	result := resultBuffer.String()
	return &result, nil
}

type DeviceService struct {
	Description string `json:"description"`
	Loaded      string `json:"loaded"`
	ServiceName string `json:"serviceName"`
	State       string `json:"state"`
	Status      string `json:"status"`
}

type DeviceServices struct {
	Running []string `json:"running"`
	Exited  []string `json:"exited"`
	Failed  []string `json:"failed"`
}

func fetchDeviceServices(name string, ctx context.Context) (*DeviceServices, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	url := fmt.Sprintf("http://%s:3333/api/services", name)
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

	var result []DeviceService
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		logger.Errorln("Fetch error:", err)
		return nil, err
	}

	var deviceServices DeviceServices

	for _, service := range result {
		if service.Status == "running" {
			deviceServices.Running = append(deviceServices.Running, service.ServiceName)
		} else if service.Status == "failed" {
			deviceServices.Failed = append(deviceServices.Failed, service.ServiceName)
		} else {
			deviceServices.Exited = append(deviceServices.Exited, service.ServiceName)
		}
	}

	return &deviceServices, nil
}

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
	if err != nil {
		logger.Errorln("Fetch error:", err)
		return nil, err
	}

	info := &DeviceInfo{
		Time:     time.Now(),
		Response: &deviceInfoResponse,
	}

	return info, nil
}

type DeviceInformation struct {
	Info    *DeviceInfo
	InfoErr error

	Status    *string
	StatusErr error

	Services    *DeviceServices
	ServicesErr error
}

func fetchForDevice(device tailscale.Device, ctx context.Context) *DeviceInformation {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var res DeviceInformation
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := fetchDeviceInfo(device.Name, ctx)
		res.Info = info
		res.InfoErr = err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		status, err := fetchDeviceStatus(device.Name, ctx)
		res.Status = status
		res.StatusErr = err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		services, err := fetchDeviceServices(device.Name, ctx)
		res.Services = services
		res.ServicesErr = err
	}()

	wg.Wait()

	return &res
}

const fetchDevicesInfoKey = "fetchDevicesInfoKey"

func fetchDevicesInfo(devices []tailscale.Device, ctx context.Context) []*DeviceInformation {

	res, found := DevicesCache.Get(fetchDevicesInfoKey)
	if found {
		return res.([]*DeviceInformation)
	}

	results := make([]*DeviceInformation, len(devices))
	var wg sync.WaitGroup

	for index, value := range devices {
		wg.Add(1)
		go func(index int, device tailscale.Device) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()

			info := fetchForDevice(device, ctx)

			results[index] = info
		}(index, value)
	}
	wg.Wait()

	DevicesCache.Set(fetchDevicesInfoKey, results, time.Minute)
	return results
}

func devices(ctx context.Context) (*Response, error) {
	filter := "tag:nixos"
	devices, err := tailscalehelper.Devices(ctx, &filter)

	if err != nil {
		return nil, err
	}

	sinceDate := time.Now()

	lastSeenDeadline := sinceDate.Add(time.Minute * -15)

	if err != nil {
		return nil, err
	}

	slices.SortFunc(devices, func(a, b tailscale.Device) int {
		return int(b.LastSeen.Time.Sub(a.LastSeen.Time).Seconds())
	})

	results := fetchDevicesInfo(devices, ctx)

	var response Response

	for i, td := range devices {
		pingResult := results[i]

		var device Device
		// timeAgo := sinceDate.Sub(td.LastSeen.Time)

		device.Id = fmt.Sprintf("device-%s", td.ID)
		device.Name = td.Hostname
		device.Ip = mainIpv4Address(td)

		if td.LastSeen.After(lastSeenDeadline) {
			device.Status = "up"
		} else {
			device.Status = "down"
		}

		device.LastSeen = td.LastSeen.Time
		device.LastSeenAgo = timeago.English.Format(td.LastSeen.Time)
		device.Adresses = td.Addresses

		if len(td.Tags) > 1 {
			device.Tags = fmt.Sprintf("%s", td.Tags)
		} else {
			device.Tags = ""
		}

		device.Hostname = td.Name
		device.Services = pingResult.Services

		var errorString strings.Builder

		err = pingResult.InfoErr
		if err != nil {
			fmt.Fprintf(&errorString, "Info: %v\n", err)
		}

		err = pingResult.ServicesErr
		if err != nil {
			fmt.Fprintf(&errorString, "Services: %v\n", err)
		}

		err = pingResult.StatusErr
		if err != nil {
			fmt.Fprintf(&errorString, "Status: %v\n", err)
		}

		device.Err = errorString.String()

		response.Devices = append(response.Devices, device)

		// timeAgo := device.LastSeen.Sub(sinceDate)
		// fmt.Fprintf(&buffer, "<pre>%d %s %s %s %s</pre>", i, device.Hostname, device.Addresses, device.LastSeen, timeAgo)
	}

	return &response, nil
}
