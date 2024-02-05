package main

import (
	_ "embed"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/tomasharkema/go-nixos-menu/server"

	"github.com/google/logger"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
)

func createLogger() {
	logPath := os.DevNull // "/tmp/nixos-devices-server.log"
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	logger.Init("nixos devices server", *verbose, true, lf)
}

func main() {
	kingpin.Parse()

	createLogger()
	logger.Infoln("start")
	server.StartServer()
}
