package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/tomasharkema/go-nixos-menu/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/logger"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	botToken = kingpin.Flag("telegram-bot-token", "Telegram bot token").Envar("TELEGRAM_BOT_TOKEN").Required().String()
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

	go startBot()

	server.StartServer()
}

func startBot() {
	bot, err := tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	hn, _ := os.Hostname()
	msg := tgbotapi.NewMessage(562728787, fmt.Sprintf("%s present!", hn))
	bot.Send(msg)

	for update := range updates {
		if update.Message != nil {
			if strings.Contains(update.Message.Text, "update") { // If we got a message
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

				text := fmt.Sprintf("Fetching devices... %s", hn)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ReplyToMessageID = update.Message.MessageID

				newMessage, err := bot.Send(msg)
				cancel()
				if err != nil {
					logger.Errorln("ERROR", hn, err)
					continue
				}

				devices, err := server.Devices(ctx)

				newText := ""
				if err != nil {
					newText = fmt.Sprintf("Error: %v %s", err, hn)
				} else {
					var msg bytes.Buffer
					for _, device := range devices.Devices {
						fmt.Fprintf(&msg, "%s: %s\n", device.Name, device.LastSeenAgo)
					}
					newText = msg.String()
					// newText = fmt.Sprintf("Devices: %v %s", devices,hn)
				}

				edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, newMessage.MessageID, newText)
				_, err = bot.Send(edit)
				if err != nil {
					logger.Errorln("ERROR", hn, err)
					continue
				}
			}
		}
	}
}
