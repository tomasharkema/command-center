package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/logger"
)

func (b *CommandBot) sendMessage(chatID int64, text string, replyMsgID *int, replyMarkup interface{}, ctx context.Context) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTMl"

	if replyMsgID != nil {
		msg.ReplyToMessageID = *replyMsgID
	}
	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	newMessage, err := b.tgBot.Send(msg)
	if err != nil {
		logger.Errorln("Error with message", err)
		return nil, err
	}

	b.storeMessageID(newMessage.MessageID)

	return &newMessage, nil
}

// func (b *CommandBot) sendErrorMessage(err error) {}

func (b *CommandBot) sendStartMessage(ctx context.Context) (*tgbotapi.Message, error) {
	welcomeString := fmt.Sprintf("<b>%s</b> present!", b.hostName)

	return b.sendMessage(b.chatID, welcomeString, nil, homeKeyboards, ctx)
}

func (b *CommandBot) sendEndMessage(ctx context.Context) (*tgbotapi.Message, error) {
	text := fmt.Sprintf("<b>%s</b> bye...", b.hostName)

	return b.sendMessage(b.chatID, text, nil, nil, ctx)
}

func (b *CommandBot) storeMessageID(msg int) {
	b.messages = append(b.messages, msg)
}
