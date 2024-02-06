package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"

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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	kingpin.Parse()

	createLogger()
	logger.Infoln("start")

	if *runBot {
		bot, err := bot.NewBot(*botToken, *chatID, ctx)
		if err != nil {
			logger.Fatalln("Error", err)
		}

		go bot.Start()
	}

	go server.StartServer(ctx)

	<-ctx.Done()
	logger.Infoln("bye")
}
