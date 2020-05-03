package users

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/storage"
)

type (
	Users interface {
		/*
			Creates user */
		CreateUser(ctx context.Context, userId, chatId uint64) error
		/*
			Registers user as a merchant */
		RegisterAsMerchant(ctx context.Context, userId uint64, merchantId, secretKey string) error

		/*
			Set active editing order for user */
		SetActiveOrderForUser(ctx context.Context, userId, orderId uint64) error
	}

	users struct {
		storage storage.Storage
	}
)

func NewUsers(storage storage.Storage) Users {
	return users{storage: storage}
}

func (u users) CreateUser(ctx context.Context, userId, chatId uint64) error {
	return u.storage.CreateUser(ctx, userId, chatId)
}

func (u users) RegisterAsMerchant(ctx context.Context, userId uint64, merchantId, secretKey string) error {
	return u.storage.SetMerchantData(ctx, userId, merchantId, secretKey)
}

func (u users) SetActiveOrderForUser(ctx context.Context, userId, orderId uint64) error {
	return u.storage.SetActiveOrderForUser(ctx, userId, orderId)
}
