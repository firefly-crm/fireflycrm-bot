package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type (
	Storage interface {
		GetOrder(ctx context.Context, orderId uint64) (types.Order, error)
		CreateOrder(ctx context.Context, userId uint64) (uint64, error)
		CreateUser(ctx context.Context, userId uint64) error
		SetMerchantData(ctx context.Context, userId uint64, merchantId, secretKey string) error
		GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error)
		SetActiveOrderMessageForUser(ctx context.Context, userId, orderId uint64) error
		GetActiveOrderForUser(ctx context.Context, userId uint64) (types.Order, error)
		UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		AddOrderMessage(ctx context.Context, orderId, messageId uint64) error
		UpdateOrderEditState(ctx context.Context, orderId uint64, state types.EditState) error
		UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error
		AddItemToOrder(ctx context.Context, orderId uint64) (receiptItemId uint64, err error)
		RemoveReceiptItem(ctx context.Context, receiptItemId uint64) error
		UpdateReceiptItemName(ctx context.Context, name string, userId, receiptItemId uint64) (err error)
		UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error)
		UpdateReceiptItemQty(ctx context.Context, qty int, receiptItemId uint64) (err error)
		GetReceiptItem(ctx context.Context, id uint64) (types.ReceiptItem, error)
		SetActiveItemId(ctx context.Context, orderId uint64, receiptItemId uint64) error
		UpdateCustomerEmail(ctx context.Context, email string, orderId uint64) (uint64, error)
		GetCustomer(ctx context.Context, customerId uint64) (c types.Customer, err error)
		AddPayment(ctx context.Context, orderId uint64, method types.PaymentMethod) (uint64, error)
		RemovePayment(ctx context.Context, paymentId uint64) error
		UpdatePaymentAmount(ctx context.Context, paymentId uint64, amount uint32) error
		UpdatePaymentLink(ctx context.Context, paymentId uint64, url, bankPaymentId string) error
		RefundPayment(ctx context.Context, paymentId uint64, amount uint32) error
		GetPayment(ctx context.Context, paymentId uint64) (types.Payment, error)
		SetActivePaymentId(ctx context.Context, orderId, paymentId uint64) error
		GetActiveOrderMessageIdForUser(ctx context.Context, userId uint64) (uint64, error)
		GetOrderMessage(ctx context.Context, messageId uint64) (types.OrderMessage, error)
		UpdateOrderMessageDisplayMode(ctx context.Context, messageId uint64, mode types.DisplayMode) error
		UpdateCustomerInstagram(ctx context.Context, instagram string, orderId uint64) (uint64, error)
		UpdateCustomerPhone(ctx context.Context, phone string, orderId uint64) (uint64, error)
		GetUser(ctx context.Context, userId uint64) (user types.User, err error)
		GetBankPayments(ctx context.Context) (payments []types.Payment, err error)
		SetPaymentPaid(ctx context.Context, paymentId uint64) error
		SetPaymentExpired(ctx context.Context, paymentId uint64) error
		GetMessagesForOrder(ctx context.Context, orderId uint64) ([]types.OrderMessage, error)
	}

	storage struct {
		db *sqlx.DB
	}
)

var (
	ErrNoSuchUser = errors.New("no such user")
)

func (s storage) GetMessagesForOrder(ctx context.Context, orderId uint64) (messages []types.OrderMessage, err error) {
	const getQuery = `SELECT * FROM order_messages WHERE order_id=$1`
	err = s.db.Select(&messages, getQuery, orderId)
	if err != nil {
		return messages, fmt.Errorf("failed to get messages: %w", err)
	}
	return
}

func (s storage) SetPaymentPaid(ctx context.Context, paymentId uint64) error {
	const updateRequest = `
UPDATE payments
SET
	payed=TRUE,
	payed_at=CURRENT_TIMESTAMP
WHERE
	id=$1`
	_, err := s.db.Exec(updateRequest, paymentId)
	if err != nil {
		return fmt.Errorf("failed to set payment paid: %w", err)
	}
	return nil
}

func (s storage) SetPaymentExpired(ctx context.Context, paymentId uint64) error {
	const updateRequest = `
UPDATE payments
SET
	expired=TRUE
WHERE
	id=$1`
	_, err := s.db.Exec(updateRequest, paymentId)
	if err != nil {
		return fmt.Errorf("failed to set payment expired: %w", err)
	}
	return nil
}

func (s storage) GetBankPayments(ctx context.Context) (payments []types.Payment, err error) {
	const getQuery = `
SELECT 
	* 
FROM
	payments
WHERE
	payment_method=0
AND
	NOT expired
AND
	NOT payed`

	err = s.db.Select(&payments, getQuery)
	if err != nil {
		return payments, fmt.Errorf("failed to get payments: %w", err)
	}
	return
}

func (s storage) GetUser(ctx context.Context, userId uint64) (user types.User, err error) {
	const getUserQuery = `SELECT * FROM users WHERE id=$1`
	err = s.db.Get(&user, getUserQuery, userId)
	if err != nil {
		return user, fmt.Errorf("failed to get user: %w", err)
	}
	return
}

func (s storage) GetPayment(ctx context.Context, paymentId uint64) (payment types.Payment, err error) {
	const getQuery = `SELECT * FROM payments WHERE id=$1`
	err = s.db.Get(&payment, getQuery, paymentId)
	if err != nil {
		return payment, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

func (s storage) UpdateCustomerInstagram(ctx context.Context, instagram string, orderId uint64) (customerId uint64, err error) {
	const createOrGetCustomerQuery = `
INSERT INTO customers(instagram) VALUES($1)
ON CONFLICT(instagram) DO UPDATE SET instagram=$1
RETURNING id
`
	const updateCustomerIdQuery = `
UPDATE orders SET customer_id=$2 WHERE id=$1
`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return customerId, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = tx.Get(&customerId, createOrGetCustomerQuery, instagram)
	if err != nil {
		return customerId, fmt.Errorf("failed to update customer instagram: %w", err)
	}

	_, err = tx.Exec(updateCustomerIdQuery, orderId, customerId)
	if err != nil {
		return customerId, fmt.Errorf("failed to update customer instagram: %w", err)
	}

	return customerId, nil
}

func (s storage) UpdateCustomerPhone(ctx context.Context, phone string, orderId uint64) (customerId uint64, err error) {
	const createOrGetCustomerQuery = `
INSERT INTO customers(phone) VALUES($1)
ON CONFLICT(phone) DO UPDATE SET phone=$1
RETURNING id
`
	const updateCustomerIdQuery = `
UPDATE orders SET customer_id=$2 WHERE id=$1
`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return customerId, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = tx.Get(&customerId, createOrGetCustomerQuery, phone)
	if err != nil {
		return customerId, fmt.Errorf("failed to update customer phone: %w", err)
	}

	_, err = tx.Exec(updateCustomerIdQuery, orderId, customerId)
	if err != nil {
		return customerId, fmt.Errorf("failed to update customer phone: %w", err)
	}

	return customerId, nil
}

func NewStorage(db *sqlx.DB) Storage {
	return storage{db}
}

func (s storage) UpdateOrderMessageDisplayMode(ctx context.Context, messageId uint64, mode types.DisplayMode) error {
	const updateQuery = `UPDATE order_messages SET display_mode=$2 WHERE id=$1`
	_, err := s.db.Exec(updateQuery, messageId, mode)
	if err != nil {
		return fmt.Errorf("failed to update display mode: %w", err)
	}
	return nil
}

func (s storage) GetOrderMessage(ctx context.Context, messageId uint64) (msg types.OrderMessage, err error) {
	const getMessageQuery = `SELECT * FROM order_messages WHERE id=$1`
	err = s.db.Get(&msg, getMessageQuery, messageId)
	if err != nil {
		return msg, fmt.Errorf("failed to get order message: %w", err)
	}
	return
}

func (s storage) GetActiveOrderMessageIdForUser(ctx context.Context, userId uint64) (msgId uint64, err error) {
	const activeMsgIdQuery = `SELECT active_order_msg_id FROM users WHERE id=$1`
	err = s.db.Get(&msgId, activeMsgIdQuery, userId)
	if err != nil {
		return msgId, fmt.Errorf("failed to get active order message id: %w", err)
	}

	return
}

func (s storage) RefundPayment(ctx context.Context, paymentId uint64, amount uint32) error {
	const updateRefundQuery = `UPDATE payments SET refund_amount=$2 WHERE id=$1`
	const updateActivePayment = `UPDATE orders SET active_payment_id=NULL WHERE id=(SELECT order_id FROM payments WHERE id=$1)`

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

	_, err = tx.Exec(updateRefundQuery, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to update refund amount: %w", err)
	}

	_, err = tx.Exec(updateActivePayment, paymentId)
	if err != nil {
		return fmt.Errorf("failed to updat active payment: %w", err)
	}
	return nil
}

func (s storage) SetActivePaymentId(ctx context.Context, orderId, paymentId uint64) error {
	const updateQuery = `UPDATE orders SET active_payment_id=$2 WHERE id=$1`
	_, err := s.db.Exec(updateQuery, orderId, paymentId)
	if err != nil {
		return fmt.Errorf("faield to set active payment id: %w", err)
	}
	return nil
}

func (s storage) UpdatePaymentLink(ctx context.Context, paymentId uint64, url, bankPaymentId string) error {
	const updateLinkQuery = `UPDATE payments SET payment_link=$2,bank_payment_id=$3 WHERE id=$1`
	_, err := s.db.Exec(updateLinkQuery, paymentId, url, bankPaymentId)
	if err != nil {
		return fmt.Errorf("failed to update payment link: %w", err)
	}
	return nil
}

func (s storage) UpdatePaymentAmount(ctx context.Context, paymentId uint64, amount uint32) error {
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

	const updateAmountQuery = `UPDATE payments SET amount=$2 WHERE id=$1`
	const updateActivePaymentQuery = `UPDATE orders SET active_payment_id=NULL WHERE id=(SELECT order_id FROM payments WHERE id=$1)`

	_, err = tx.Exec(updateAmountQuery, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to update payment amount: %w", err)
	}

	_, err = tx.Exec(updateActivePaymentQuery, paymentId)
	if err != nil {
		return fmt.Errorf("failed to update active payment id: %w", err)
	}

	return nil
}

func (s storage) RemovePayment(ctx context.Context, paymentId uint64) error {
	const removeQuery = `DELETE FROM payments WHERE id=$1`
	_, err := s.db.Exec(removeQuery, paymentId)
	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}
	return nil
}

func (s storage) AddPayment(ctx context.Context, orderId uint64, method types.PaymentMethod) (id uint64, err error) {
	const addPaymentQuery = `
INSERT INTO payments(order_id, payment_method, payed, payed_at)
VALUES($1, $2, $3, $4)
RETURNING id`

	const setActivePaymentQuery = `
UPDATE orders SET active_payment_id=$2 WHERE id=$1
`

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

	payed := false
	var payedAt sql.NullTime
	if method == types.PaymentMethodCash || method == types.PaymentMethodCard2Card {
		payed = true
		payedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	err = tx.Get(&id, addPaymentQuery, orderId, method, payed, payedAt)
	if err != nil {
		return id, fmt.Errorf("failed to add payment: %w", err)
	}

	_, err = tx.Exec(setActivePaymentQuery, orderId, id)
	if err != nil {
		return id, fmt.Errorf("failed to set active bill id: %w", err)
	}

	return
}

func (s storage) GetCustomer(ctx context.Context, customerId uint64) (c types.Customer, err error) {
	const getCustomerQuery = `SELECT * FROM customers WHERE id=$1`
	err = s.db.Get(&c, getCustomerQuery, customerId)
	if err != nil {
		return c, fmt.Errorf("failed to get customer: %w", err)
	}

	fmt.Printf("CUSTOMER: %v\n", c)

	return c, nil
}

func (s storage) UpdateCustomerEmail(ctx context.Context, email string, orderId uint64) (customerId uint64, err error) {
	const createOrGetCustomerQuery = `
INSERT INTO customers(email) VALUES($1)
ON CONFLICT(email) DO UPDATE SET email=$1
RETURNING id
`
	const updateCustomerIdQuery = `
UPDATE orders SET customer_id=$2 WHERE id=$1
`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return customerId, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = tx.Get(&customerId, createOrGetCustomerQuery, email)
	if err != nil {
		return customerId, fmt.Errorf("failed to update customer email: %w", err)
	}

	_, err = tx.Exec(updateCustomerIdQuery, orderId, customerId)
	if err != nil {
		return customerId, fmt.Errorf("failed to update customer email: %w", err)
	}

	return customerId, nil
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
	const getOrderQuery = `SELECT * FROM orders WHERE id=$1`
	const getReceiptItemsQuery = `SELECT * FROM receipt_items WHERE order_id=$1 `
	const getPaymentsQuery = `SELECT * FROM payments WHERE order_id=$1`

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

	err = tx.Get(&order, getOrderQuery, orderId)
	if err != nil {
		return order, fmt.Errorf("failed to get order: %w", err)
	}

	err = tx.Select(&order.ReceiptItems, getReceiptItemsQuery, orderId)
	if err != nil {
		return order, fmt.Errorf("failed to get receipt items: %w", err)
	}

	err = tx.Select(&order.Payments, getPaymentsQuery, orderId)
	if err != nil {
		return order, fmt.Errorf("failed to get payments: %w", err)
	}

	return
}

func (s storage) GetReceiptItem(ctx context.Context, id uint64) (item types.ReceiptItem, err error) {
	const getItemQuery = `SELECT * FROM receipt_items WHERE id=$1`
	err = s.db.Get(&item, getItemQuery, id)
	if err != nil {
		return item, fmt.Errorf("failed to get receipt item: %w", err)
	}

	return
}

func (s storage) UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error) {
	const updateQuery = `UPDATE receipt_items SET price=$2,initialised=TRUE WHERE id=$1`
	const updateActiveItem = `UPDATE orders SET active_item_id=NULL WHERE id=(SELECT order_id FROM receipt_items WHERE id=$1)`

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

	_, err = tx.Exec(updateQuery, receiptItemId, price)
	if err != nil {
		return fmt.Errorf("failed to update item price: %w", err)
	}

	_, err = tx.Exec(updateActiveItem, receiptItemId)
	if err != nil {
		return fmt.Errorf("failed to update active item id: %w", err)
	}
	return nil
}

func (s storage) UpdateReceiptItemQty(ctx context.Context, qty int, receiptItemId uint64) (err error) {
	const updateQuery = `UPDATE receipt_items SET quantity=$2 WHERE id=$1`
	_, err = s.db.Exec(updateQuery, receiptItemId, qty)
	if err != nil {
		return fmt.Errorf("failed to update item quantity: %w", err)
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

func (s storage) SetActiveItemId(ctx context.Context, orderId uint64, receiptItemId uint64) error {
	const setActiveItemIdQuery = `UPDATE orders SET active_item_id=$1 WHERE id=$2`
	_, err := s.db.Exec(setActiveItemIdQuery, receiptItemId, orderId)
	if err != nil {
		return fmt.Errorf("failed to set active item id: %w", err)
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

func (s storage) UpdateOrderEditState(ctx context.Context, orderId uint64, state types.EditState) error {
	const updateQuery = `UPDATE orders SET edit_state=$1 WHERE id=$2`
	_, err := s.db.Exec(updateQuery, state, orderId)
	if err != nil {
		return fmt.Errorf("failed to update order edit state: %w", err)
	}
	return nil
}

func (s storage) UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error {
	const updateQuery = `UPDATE orders SET order_state=$1 WHERE id=$2`
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
	* 
FROM 
	orders 
WHERE 
	id=(SELECT order_id FROM order_messages WHERE id=(
	    SELECT active_order_msg_id FROM users WHERE id=$1))`
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

func (s storage) AddOrderMessage(ctx context.Context, orderId, messageId uint64) error {
	const addQuery = `INSERT INTO order_messages(id, order_id) VALUES($1,$2)`
	_, err := s.db.Exec(addQuery, messageId, orderId)
	if err != nil {
		return fmt.Errorf("failed to update message id: %w", err)
	}
	return nil
}

func (s storage) SetActiveOrderMessageForUser(ctx context.Context, userId, messageId uint64) error {
	log.Printf("seting active order; userId: %d, orderId: %d", userId, messageId)
	const setActiveOrderQuery = `UPDATE users SET active_order_msg_id=$1 WHERE id=$2`
	_, err := s.db.Exec(setActiveOrderQuery, messageId, userId)
	if err != nil {
		return fmt.Errorf("failed to set active order: %w", err)
	}
	return nil
}

func (s storage) GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error) {
	const getOrderIdQuery = `SELECT order_id FROM order_messages WHERE id=$1`

	const getOrderQuery = `SELECT * FROM orders WHERE id=$1`

	const getReceiptItemsQuery = `SELECT * FROM receipt_items WHERE order_id=$1
`
	const getPaymentsQuery = `SELECT * FROM payments WHERE order_id=$1`

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

	var orderId uint64
	err = tx.Get(&orderId, getOrderIdQuery, messageId)
	if err != nil {
		return order, fmt.Errorf("failed to get order id from messages: %w", err)
	}

	err = tx.Get(&order, getOrderQuery, orderId)
	if err != nil {
		return order, fmt.Errorf("failed to get order: %w", err)
	}

	err = tx.Select(&order.ReceiptItems, getReceiptItemsQuery, order.Id)
	if err != nil {
		return order, fmt.Errorf("failed to get receipt items: %w", err)
	}

	err = tx.Select(&order.Payments, getPaymentsQuery, order.Id)
	if err != nil {
		return order, fmt.Errorf("failed to get payments: %w", err)
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

func (s storage) CreateUser(ctx context.Context, userId uint64) error {
	const checkUserExists = `SELECT EXISTS(SELECT * FROM users WHERE id=$1)`
	const createUserQuery = `INSERT INTO users (id) VALUES ($1)`

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

	_, err = tx.Exec(createUserQuery, userId)
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
