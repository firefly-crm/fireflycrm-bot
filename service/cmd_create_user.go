package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) createUser(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	userId := uint64(update.Message.From.ID)
	err := s.Users.CreateUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	msg := tg.NewMessage(update.Message.Chat.ID, replyWelcome)
	msg.ReplyMarkup = merchantStandByKeyboardMarkup()
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w")
	}

	return nil
}
