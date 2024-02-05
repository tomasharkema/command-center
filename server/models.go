package server

import "time"

type Device struct {
	Id       string
	Name     string
	Status   string
	LastSeen string
	Adresses string
	Tags     string
	Hostname string
	Ip       string
	Err      *string
}

type Response struct {
	Devices []Device
}

type DeviceInfo struct {
	Time     time.Time
	Response *DeviceInfoResponse
	Err      *error
}

type DeviceInfoResponse struct {
	Label string
	Name  string

	Revision string
	Tags     []string
	Version  string
}
