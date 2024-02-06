package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/logger"
	"github.com/tomasharkema/command-center/server"
)

type CommandBot struct {
	tgBot  *tgbotapi.BotAPI
	chatID int64
	ctx    context.Context
}

var homeKeyboards = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Fetch newest update", "/update"),
	),
	// tgbotapi.NewInlineKeyboardRow(
	// 	tgbotapi.NewInlineKeyboardButtonData("/clean", "/clean"),
	// ),
)

func NewBot(botToken string, chatID int64, ctx context.Context) (*CommandBot, error) {
	cb := &CommandBot{chatID: chatID, ctx: ctx}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	cb.tgBot = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return cb, nil
}

func (b *CommandBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tgBot.GetUpdatesChan(u)

	go b.sendStartMessage()

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				b.handleCommand(update.Message.Chat.ID, update.Message.Text)
			} else {
				b.handleText(update.Message)
			}
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		} else {
			logger.Warningf("no clue %v", update)
		}

	}

}

func (b *CommandBot) handleCommand(chatID int64, txt string) {
	fmt.Println("HANDLE COMMAND!", txt)
	switch txt {
	case "/update":
		b.handleUpdateCommand(chatID)
	}
}

func (b *CommandBot) handleText(msg *tgbotapi.Message) {
	// switch msg.Text {
	// case "/update":
	// 	b.handleUpdateCommand(msg)
	// }
}

func (b *CommandBot) handleCallback(callbackInfo *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(callbackInfo.ID, callbackInfo.Data)
	b.handleCommand(callbackInfo.Message.Chat.ID, callback.Text)

	if _, err := b.tgBot.Request(callback); err != nil {
		logger.Errorln("Error", err)
	}

	// msg := tgbotapi.NewMessage(callbackInfo.Message.Chat.ID, callbackInfo.Data)
	// if _, err := b.tgBot.Send(msg); err != nil {
	// 	logger.Errorln("Error", err)
	// }
}

func (b *CommandBot) sendStartMessage() {
	hn, _ := os.Hostname()

	welcomeString := fmt.Sprintf("<b>%s</b> present!", hn)

	msg := tgbotapi.NewMessage(b.chatID, welcomeString)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = homeKeyboards

	b.tgBot.Send(msg)
}

func (b *CommandBot) handleUpdateCommand(chatID int64) {

	text := fmt.Sprintln("<i>Fetching devices...</i>")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTMl"
	// msg.ReplyToMessageID = message.MessageID
	newMessage, err := b.tgBot.Send(msg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	devices, err := server.Devices(ctx)
	newText := ""
	var buttons [][]tgbotapi.InlineKeyboardButton

	if err != nil {
		newText = fmt.Sprintf("Error: %v", err)
	} else {

		for _, device := range devices.Devices {
			text := fmt.Sprintf("🟢 %s", device.Name)

			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(text, "/update"),
			)
			buttons = append(buttons, row)
		}
		newText = "<b>Got the following devices:</b>"
		// newText = fmt.Sprintf("Devices: %v %s", devices,hn)
	}

	rows := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	edit := tgbotapi.NewEditMessageTextAndMarkup(b.chatID, newMessage.MessageID, newText, rows)
	edit.ParseMode = "HTML"
	_, err = b.tgBot.Send(edit)
	if err != nil {
		logger.Errorln("ERROR", err)
	}

}

// func startBot() {
// 	bot, err := tgbotapi.NewBotAPI(*botToken)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bot.Debug = true

// 	log.Printf("Authorized on account %s", bot.Self.UserName)

// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60

// 	updates := bot.GetUpdatesChan(u)

// 	hn, _ := os.Hostname()
// 	msg := tgbotapi.NewMessage(562728787, fmt.Sprintf("%s present!", hn))
// 	bot.Send(msg)

// 	for update := range updates {
// 		if update.Message != nil {
// 			if strings.Contains(update.Message.Text, "update") { // If we got a message
// 				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// 				text := fmt.Sprintf("Fetching devices... %s", hn)

// 				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
// 				msg.ReplyToMessageID = update.Message.MessageID

// 				newMessage, err := bot.Send(msg)
// 				cancel()
// 				if err != nil {
// 					logger.Errorln("ERROR", hn, err)
// 					continue
// 				}

// 				devices, err := server.Devices(ctx)

// 				newText := ""
// 				if err != nil {
// 					newText = fmt.Sprintf("Error: %v %s", err, hn)
// 				} else {
// 					var msg bytes.Buffer
// 					for _, device := range devices.Devices {
// 						fmt.Fprintf(&msg, "%s: %s\n", device.Name, device.LastSeenAgo)
// 					}
// 					newText = msg.String()
// 					// newText = fmt.Sprintf("Devices: %v %s", devices,hn)
// 				}

// 				edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, newMessage.MessageID, newText)
// 				_, err = bot.Send(edit)
// 				if err != nil {
// 					logger.Errorln("ERROR", hn, err)
// 					continue
// 				}
// 			}
// 		}
// 	}
// }
