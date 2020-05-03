package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/jmoiron/sqlx"
	"log"
)

type (
	Storage interface {
		GetOrder(ctx context.Context, orderId uint64) (types.Order, error)
		CreateOrder(ctx context.Context, userId uint64) (uint64, error)
		CreateUser(ctx context.Context, userId, chatId uint64) error
		SetMerchantData(ctx context.Context, userId uint64, merchantId, secretKey string) error
		GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error)
		SetActiveOrderForUser(ctx context.Context, userId, orderId uint64) error
		GetActiveOrderForUser(ctx context.Context, userId uint64) (types.Order, error)
		UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		UpdateMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error
		AddItemToOrder(ctx context.Context, orderId uint64) (receiptItemId uint64, err error)
		RemoveReceiptItem(ctx context.Context, receiptItemId uint64) error
		UpdateReceiptItemName(ctx context.Context, name string, userId, receiptItemId uint64) (err error)
		UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error)
		GetReceiptItem(ctx context.Context, id uint64) (types.ReceiptItem, error)
		GetReceiptItems(ctx context.Context, id uint64) ([]types.ReceiptItem, error)
	}

	storage struct {
		db *sqlx.DB
	}
)

func (s storage) GetReceiptItems(ctx context.Context, id uint64) ([]types.ReceiptItem, error) {
	//TODO: Implement
	panic("implement me")
}

var (
	ErrNoSuchUser = errors.New("no such user")
)

func NewStorage(db *sqlx.DB) Storage {
	return storage{db}
}

func (s storage) RemoveReceiptItem(ctx context.Context, receiptItemId uint64) error {
	const deleteQuery = `DELETE FROM receipt_items WHERE id=$1`
	_, err := s.db.Exec(deleteQuery, receiptItemId)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (s storage) GetOrder(ctx context.Context, orderId uint64) (order types.Order, err error) {
	const getOrderQuery = `
SELECT 
	id,
	message_id,
	user_id,
	description,
	state,
    active_item_id,
    hint_message_id
FROM
	orders
WHERE
	id=$1`

	const getReceiptItemsQuery = `
SELECT
	id,
    name,
    item_id,
    order_id,
    quantity,
    price,
    initialised
FROM
	receipt_items
WHERE
	order_id=$1
`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return order, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = s.db.Get(&order, getOrderQuery, orderId)
	if err != nil {
		return order, fmt.Errorf("failed to get order: %w", err)
	}

	err = s.db.Select(&order.ReceiptItems, getReceiptItemsQuery, orderId)
	if err != nil {
		return order, fmt.Errorf("failed to get receipt items: %w", err)
	}

	return
}

func (s storage) GetReceiptItem(ctx context.Context, id uint64) (item types.ReceiptItem, err error) {
	const getItemQuery = `
SELECT
	id,
	name,
	order_id,
	item_id,
	price,
	quantity,
	initialised
FROM
	receipt_items
WHERE
	id=$1
`
	err = s.db.Get(&item, getItemQuery, id)
	if err != nil {
		return item, fmt.Errorf("failed to get receipt item: %w", err)
	}

	return
}

func (s storage) UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error) {
	const updateQuery = `UPDATE receipt_items SET price=$2,initialised=TRUE WHERE id=$1`
	_, err = s.db.Exec(updateQuery, receiptItemId, price)
	if err != nil {
		return fmt.Errorf("failed to update item price: %w", err)
	}
	return nil
}

func (s storage) UpdateReceiptItemName(ctx context.Context, name string, userId, receiptItemId uint64) (err error) {
	const getItemQuery = `
INSERT INTO items(user_id,name)
VALUES($1,$2)
ON CONFLICT(user_id,name) DO UPDATE SET name=$2
RETURNING id
`
	const updateNameQuery = `
UPDATE receipt_items
SET
	item_id=$2,
	name=$3
WHERE
	id=$1`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var itemId uint64
	err = tx.Get(&itemId, getItemQuery, userId, name)
	if err != nil {
		return fmt.Errorf("failed to create or get item: %w", err)
	}

	_, err = tx.Exec(updateNameQuery, receiptItemId, itemId, name)
	if err != nil {
		return fmt.Errorf("failed to update receipt item: %w", err)
	}

	return nil
}

func (s storage) AddItemToOrder(ctx context.Context, orderId uint64) (receiptItemId uint64, err error) {
	const createReceiptItemQuery = `INSERT INTO receipt_items(order_id) VALUES ($1) RETURNING id`

	const setActiveItemIdQuery = `UPDATE orders SET active_item_id=$1 WHERE id=$2`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return receiptItemId, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = tx.Get(&receiptItemId, createReceiptItemQuery, orderId)
	if err != nil {
		return receiptItemId, fmt.Errorf("failed to create receipt item: %w", err)
	}

	_, err = tx.Exec(setActiveItemIdQuery, receiptItemId, orderId)
	if err != nil {
		return receiptItemId, fmt.Errorf("failed to set active item id for order: %w", err)
	}

	return
}

func (s storage) UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error {
	const updateQuery = `UPDATE orders SET state=$1 WHERE id=$2`
	_, err := s.db.Exec(updateQuery, state, orderId)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}
	return nil
}

func (s storage) GetActiveOrderForUser(ctx context.Context, userId uint64) (types.Order, error) {
	var order types.Order

	const getOrderQuery = `
SELECT 
	id,
	message_id,
	user_id,
	description,
	state,
    active_item_id,
    hint_message_id
FROM
	orders
WHERE
	id=(SELECT active_order_id FROM users WHERE id=$1)
`
	err := s.db.Get(&order, getOrderQuery, userId)
	if err != nil {
		return order, fmt.Errorf("failed to get active order for user id: %w", err)
	}
	return order, nil
}

func (s storage) UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error {
	const updateQuery = `UPDATE orders SET hint_message_id=$1 WHERE id=$2`
	_, err := s.db.Exec(updateQuery, messageId, orderId)
	if err != nil {
		return fmt.Errorf("failed to update hint message id: %w", err)
	}
	return nil
}

func (s storage) UpdateMessageForOrder(ctx context.Context, orderId, messageId uint64) error {
	const updateQuery = `UPDATE orders SET message_id=$1 WHERE id=$2`
	_, err := s.db.Exec(updateQuery, messageId, orderId)
	if err != nil {
		return fmt.Errorf("failed to update message id: %w", err)
	}
	return nil
}

func (s storage) SetActiveOrderForUser(ctx context.Context, userId, orderId uint64) error {
	log.Printf("seting active order; userId: %d, orderId: %d", userId, orderId)
	const setActiveOrderQuery = `UPDATE users SET active_order_id=$1 WHERE id=$2`
	_, err := s.db.Exec(setActiveOrderQuery, orderId, userId)
	if err != nil {
		return fmt.Errorf("failed to set active order: %w", err)
	}
	return nil
}

func (s storage) GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error) {
	const getOrderQuery = `SELECT id,message_id,user_id,description,state,active_item_id,hint_message_id FROM orders WHERE message_id=$1`
	err = s.db.Get(&order, getOrderQuery, messageId)
	if err != nil {
		return order, fmt.Errorf("failed to get order by message id: %w", err)
	}

	return
}

func (s storage) SetMerchantData(ctx context.Context, userId uint64, merchantId, secretKey string) error {
	const updateDataQuery = `
UPDATE 
	users 
SET
	is_merchant=TRUE,
	merchant_id=$2,
	secret_key=$3
WHERE
	id=$1`

	_, err := s.db.Exec(updateDataQuery, userId, merchantId, secretKey)
	if err != nil {
		return fmt.Errorf("failed to set user as merchant: %w", err)
	}

	return nil
}

func (s storage) CreateUser(ctx context.Context, userId, chatId uint64) error {
	const checkUserExists = `SELECT EXISTS(SELECT * FROM users WHERE id=$1)`
	const createUserQuery = `INSERT INTO users (id,chat_id) VALUES ($1,$2)`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start create user transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	isExists := false

	err = tx.Get(&isExists, checkUserExists, userId)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	//If user already exists it's ok-flow
	if isExists {
		return nil
	}

	_, err = tx.Exec(createUserQuery, userId, chatId)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s storage) CreateOrder(ctx context.Context, userId uint64) (id uint64, err error) {
	const checkUserExists = `SELECT EXISTS(SELECT * FROM users WHERE id=$1)`
	const createOrderQuery = `INSERT INTO orders(user_id) VALUES ($1) RETURNING id`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return id, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	isExists := false

	err = tx.Get(&isExists, checkUserExists, userId)
	if err != nil {
		return id, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if !isExists {
		return id, ErrNoSuchUser
	}

	err = tx.Get(&id, createOrderQuery, userId)
	if err != nil {
		return id, fmt.Errorf("failed to create new order: %w", err)
	}

	return id, nil
}
