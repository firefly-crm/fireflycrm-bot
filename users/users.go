package users

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	. "github.com/DarthRamone/fireflycrm-bot/types"
)

type (
	Users interface {
		/*
			Creates user */
		CreateUser(ctx context.Context, userId uint64) error
		/*
			Registers user as a merchant */
		RegisterAsMerchant(ctx context.Context, userId uint64, merchantId, secretKey string) error

		/*
			Set active editing order for user */
		SetActiveOrderMessageForUser(ctx context.Context, userId, orderId uint64) error

		GetCustomer(ctx context.Context, customerId uint64) (c Customer, err error)
	}

	users struct {
		storage storage.Storage
	}
)

func NewUsers(storage storage.Storage) Users {
	return users{storage: storage}
}

func (u users) CreateUser(ctx context.Context, userId uint64) error {
	return u.storage.CreateUser(ctx, userId)
}

func (u users) RegisterAsMerchant(ctx context.Context, userId uint64, merchantId, secretKey string) error {
	return u.storage.SetMerchantData(ctx, userId, merchantId, secretKey)
}

func (u users) SetActiveOrderMessageForUser(ctx context.Context, userId, orderId uint64) error {
	return u.storage.SetActiveOrderMessageForUser(ctx, userId, orderId)
}

func (u users) GetCustomer(ctx context.Context, customerId uint64) (c Customer, err error) {
	return u.storage.GetCustomer(ctx, customerId)
}
