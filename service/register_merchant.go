package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (s Service) registerMerchant(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	cmd := update.Message.Text
	userId := uint64(update.Message.From.ID)

	args := strings.Split(cmd, " ")
	if len(args) != 3 {
		return fmt.Errorf("wrong arguments count")
	}

	merchantId := args[1]
	secretKey := args[2]

	err := s.Users.RegisterAsMerchant(ctx, userId, merchantId, secretKey)
	if err != nil {
		return fmt.Errorf("failed to register as merchant: %w", err)
	}

	deleteMessage := tg.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err = bot.DeleteMessage(deleteMessage)
	if err != nil {
		return fmt.Errorf("failed to delete command message: %w", err)
	}

	msg := tg.NewMessage(update.Message.Chat.ID, replyMerchantSuccessfulRegistered)
	msg.ReplyMarkup = merchantStandByKeyboardMarkup()
	msg.ParseMode = "markdown"

	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
