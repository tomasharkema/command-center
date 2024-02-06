package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/tomasharkema/command-center/bot"
	"github.com/tomasharkema/command-center/server"

	"github.com/google/logger"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	botToken = kingpin.Flag("telegram-bot-token", "Telegram bot token").Envar("TELEGRAM_BOT_TOKEN").Required().String()
	chatID   = kingpin.Flag("telegram-chat-id", "Telegram bot token").Envar("TELEGRAM_CHAT_ID").Required().Int64()
	runBot   = kingpin.Flag("run-bot", "Telegram bot token").Envar("COMMAND_CENTER_RUN_BOT").Default("false").Bool()

	ip = kingpin.Flag("listen", "IP address to ping.").Short('l').Default(":3333").TCP()

	commandBot *bot.CommandBot
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	kingpin.Parse()

	createLogger()
	logger.Infoln("start")

	if *runBot {
		b, err := bot.NewBot(*botToken, *chatID, ctx)
		if err != nil {
			logger.Fatalln("Error", err)
		}
		commandBot = b
		go b.Start()
	}

	go server.StartServer((*ip).String(), *runBot, ctx)

	logger.Infoln("Waiting for signal...")
	sig := <-ctx.Done()
	logger.Warningln("Closing with", sig)
	cleanups()

	logger.Infoln("bye")
}

func cleanups() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	logger.Infoln("Run cleanups...")
	if commandBot != nil {
		commandBot.Cleanups(ctx)
	}
	logger.Infoln("Run cleanups done")
}
