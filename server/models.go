package server

import "time"

type Device struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Status      string          `json:"status"`
	LastSeen    time.Time       `json:"lastSeen"`
	LastSeenAgo string          `json:"lastSeenAgo"`
	Adresses    []string        `json:"addresses"`
	Tags        string          `json:"tags"`
	Hostname    string          `json:"hostname"`
	Ip          string          `json:"ip"`
	Err         string          `json:"error"`
	Services    *DeviceServices `json:"services"`
}

type Response struct {
	Devices []Device `json:"devices"`
}

type DeviceInfo struct {
	Time     time.Time           `json:"time"`
	Response *DeviceInfoResponse `json:"response"`
	Err      *error              `json:"error"`
}

type DeviceInfoResponse struct {
	Label string `json:"label"`
	Name  string `json:"name"`

	Revision string   `json:"revision"`
	Tags     []string `json:"tags"`
	Version  string   `json:"version"`
}
