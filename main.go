package main

import (
	_ "embed"
	"os"

	"github.com/tomasharkema/go-nixos-menu/server"

	"github.com/google/logger"
)

func createLogger() {
	logPath := "/tmp/nixos-devices-server.log"
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	logger.Init("nixos devices server", true, true, lf)
}

func main() {
	createLogger()
	logger.Infoln("start")
	server.StartServer()
}
