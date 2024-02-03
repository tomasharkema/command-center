package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"sync"
	"text/template"
	"time"

	"github.com/tailscale/tailscale-client-go/tailscale"
	ts "github.com/tomasharkema/go-nixos-menu/tailscale"
	"github.com/xeonx/timeago"
)

//go:embed assets/devices.html
var devicesHtml string

type Device struct {
	Id       string
	Name     string
	Status   string
	LastSeen string
	Adresses string
	Tags     string
	Hostname string
}

type Response struct {
	Devices []Device
}

type DeviceInfo struct {
	
}

func fetchDeviceInfo(name string, ctx context.Context) (result *string, err error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second * 10)
	defer cancel()

	url := fmt.Sprintf("https://%s/webhook/hooks/info-json", name)
	log.Printf("Fetch url for: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		// results[index] =  fmt.Sprintf("ERR %s", err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// results[index] = fmt.Sprintf("ERR %s", err)
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		// results[index] = fmt.Sprintf("ERR %s", err)
		return nil, err
	}
info:=fmt.Sprintf("JAJ %s", b)
	return &info, nil

	// results[index] = fmt.Sprintf("JAJ %s", b)
}


func fetchDevicesInfo(devices []tailscale.Device, ctx context.Context) []string {
	results := make([]string, len(devices))
	var wg sync.WaitGroup

	for index, value := range devices {
		wg.Add(1)
		go func(index int, device tailscale.Device) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, time.Second * 10)
			defer cancel()

			url := fmt.Sprintf("https://%s/webhook/hooks/info-json", device.Name)
			log.Println(url)

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				results[index] =  fmt.Sprintf("ERR %s", err)
				return
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				results[index] = fmt.Sprintf("ERR %s", err)
				return
			}
			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			if err != nil {
				results[index] = fmt.Sprintf("ERR %s", err)
				return
			}
		
			results[index] = fmt.Sprintf("JAJ %s", b)
		}(index, value)
	}
	wg.Wait()
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	devices, err := ts.Devices(ctx)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	sinceDate := time.Now()

	lastSeenDeadline := sinceDate.Add(time.Minute * -15)

	t, err := template.New("foo").Parse(devicesHtml)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	results := fetchDevicesInfo(devices)

	var response Response

	for i, td := range devices {
		pingResult := results[i]
		if !slices.Contains(td.Tags, "tag:nixos") {
			continue
		}

		var device Device
		// timeAgo := sinceDate.Sub(td.LastSeen.Time)

		device.Id = fmt.Sprintf("device-%s", td.ID)
		device.Name = td.Hostname

		if td.LastSeen.After(lastSeenDeadline) {
			device.Status = "up"
		} else {
			device.Status = "down"
		}

		device.LastSeen = timeago.English.Format(td.LastSeen.Time)
		device.Adresses = fmt.Sprintf("%s", td.Addresses)

		if len(td.Tags) > 1 {
			device.Tags = fmt.Sprintf("%s", td.Tags)
		} else {
			device.Tags = ""
		}

		device.Hostname = fmt.Sprintf("%s, %s", td.Name, pingResult)

		response.Devices = append(response.Devices, device)

		// timeAgo := device.LastSeen.Sub(sinceDate)
		// fmt.Fprintf(&buffer, "<pre>%d %s %s %s %s</pre>", i, device.Hostname, device.Addresses, device.LastSeen, timeAgo)
	}

	err = t.Execute(w, response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	http.HandleFunc("/", devicesHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
