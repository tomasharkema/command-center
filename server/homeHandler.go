package server

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed assets/devices.html
var devicesHtml string

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	// defer cancel()

	response := GetDevices()
	if response == nil {
		http.Error(w, "No response", 500)
		return
	}

	t, err := template.New("foo").Parse(devicesHtml)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
