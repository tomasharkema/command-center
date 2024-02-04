package server

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	ts "github.com/tomasharkema/go-nixos-menu/tailscale"
	"github.com/xeonx/timeago"
)

//go:embed assets/devices.html
var devicesHtml string

func homeHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	filter := "tag:nixos"
	devices, err := ts.Devices(ctx, &filter)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	sinceDate := time.Now()

	lastSeenDeadline := sinceDate.Add(time.Minute * -15)

	t, err := template.New("foo").Parse(devicesHtml)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	results := fetchDevicesInfo(devices, ctx)

	var response Response

	for i, td := range devices {
		pingResult := results[i]

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