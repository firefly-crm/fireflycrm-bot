package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/jmoiron/sqlx"
)

type (
	Storage interface {
		CreateOrder(ctx context.Context, opts types.OrderOptions) (uint64, error)
	}

	storage struct {
		db *sqlx.DB
	}
)

func NewStorage(db *sqlx.DB) Storage {
	return storage{db}
}

func (s storage) CreateOrder(ctx context.Context, opts types.OrderOptions) (id uint64, err error) {
	const getCustomerQuery = `SELECT id FROM customers WHERE email=$1 OR phone=$2 OR social=$3`

	const createCustomerQuery = `
INSERT INTO customers (
    name,
    email,
    phone,
    social
) VALUES (
	$1,
	$2,
	$3,
	$4
)
RETURNING id`

	const createOrderQuery = `
INSERT INTO orders(
	customer_id,
	description
) VALUES (
	$1,
	$2
)
RETURNING id`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	customer := &types.Customer{}
	err = tx.Get(customer, getCustomerQuery, opts.CustomerEmail, opts.CustomerPhone, opts.CustomerSocial)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var customerId uint64
			err := tx.Get(&customerId, createCustomerQuery, opts.CustomerName, opts.CustomerEmail, opts.CustomerPhone, opts.CustomerSocial)
			if err != nil {
				return 0, fmt.Errorf("failed to create customer: %w", err)
			}
			customer.Id = customerId
		} else {
			return 0, fmt.Errorf("failed to get customer id: %w", err)
		}
	}

	var orderId uint64
	err = tx.Get(&orderId, createOrderQuery, customer.Id, opts.Description)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	return orderId, nil
}
