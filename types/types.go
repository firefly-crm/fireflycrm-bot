package types

import (
	"database/sql"
	"fmt"
)

const (
	GOODS ReceiptItemType = iota
	DELIVERY
)

const (
	FULL_PREPAYMENT PaymentType = iota
	PARTIAL_PREPAYMENT
)

const (
	StandBy OrderState = iota
	WaitingItemName
	WaitingItemPrice
	WaitingItemQuantity
	WaitingCustomerEmail
	WaitingCustomerPhone
	WaitingCustomerName
	Completed = 99
)

type (
	OrderState byte

	ReceiptItemType byte

	PaymentType byte

	OrderOptions struct {
		Description    string
		PaymentType    PaymentType
		CustomerName   string
		CustomerEmail  string
		CustomerPhone  string
		CustomerSocial string
	}

	ReceiptItem struct {
		Id          uint64        `db:"id"`
		Name        string        `db:"name"`
		ItemId      sql.NullInt64 `db:"item_id"`
		OrderId     uint64        `db:"order_id"`
		Price       uint32        `db:"price"`
		Quantity    uint32        `db:"quantity"`
		Initialised bool          `db:"initialised"`
	}

	Customer struct {
		Id uint64
	}

	Order struct {
		Id            uint64        `db:"id"`
		HintMessageId sql.NullInt64 `db:"hint_message_id"`
		MessageId     uint64        `db:"message_id"`
		UserId        uint64        `db:"user_id"`
		Description   string        `db:"description"`
		ActiveItemId  sql.NullInt64 `db:"active_item_id"`
		State         OrderState    `db:"state"`
		ReceiptItems  []ReceiptItem
	}

	Bill struct {
		Id  uint64
		Url string
	}
)

func (o Order) MessageString() string {
	result :=
		`*Заказ №%d*

*Позиции:*
`

	var amount float32
	if o.ReceiptItems != nil {
		for _, i := range o.ReceiptItems {
			price := float32(i.Price*i.Quantity) / 100.0
			amount += price
			result += fmt.Sprintf("- %s\t\t%.2f₽\tx%d\n", i.Name, price, i.Quantity)
		}
	}
	result += fmt.Sprintf("*Итого: %.2f₽\n*", amount)

	return fmt.Sprintf(result, o.Id)
}
