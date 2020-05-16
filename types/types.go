package types

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	PaymentMethodAcquiring PaymentMethod = iota
	PaymentMethodCard2Card
	PaymentMethodCash
)

const (
	EditStateNone EditState = iota
	EditStateWaitingItemName
	EditStateWaitingItemPrice
	EditStateWaitingItemQuantity
	EditStateWaitingCustomerEmail
	EditStateWaitingCustomerInstagram
	EditStateWaitingCustomerPhone
	EditStateWaitingPaymentAmount
	EditStateWaitingRefundAmount
)

const (
	OrderStateForming OrderState = iota
	OrderStateInProgress
	OrderStateDone
	OrderStateDeleted = 99
)

const (
	DisplayModeFull DisplayMode = iota
	DisplayModeCollapsed
	DisplayModeDeleted
)

const (
	ReceiptItemTypeGoods ReceiptItemType = iota
	ReceiptItemTypeDelivery
)

type (
	OrderState      byte
	EditState       byte
	ReceiptItemType byte
	PaymentMethod   byte
	DisplayMode     byte

	OrderOptions struct {
		Description    string
		PaymentType    PaymentMethod
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
		CreatedAt   time.Time     `db:"created_at"`
		UpdatedAt   time.Time     `db:"updated_at"`
	}

	Customer struct {
		Id        uint64         `db:"id"`
		Email     sql.NullString `db:"email"`
		Phone     sql.NullString `db:"phone"`
		Name      sql.NullString `db:"name"`
		Instagram sql.NullString `db:"instagram"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
	}

	Order struct {
		Id              uint64        `db:"id"`
		HintMessageId   sql.NullInt64 `db:"hint_message_id"`
		UserId          uint64        `db:"user_id"`
		CustomerId      sql.NullInt64 `db:"customer_id"`
		Description     string        `db:"description"`
		ActiveItemId    sql.NullInt64 `db:"active_item_id"`
		ActivePaymentId sql.NullInt64 `db:"active_payment_id"`
		OrderState      OrderState    `db:"order_state"`
		EditState       EditState     `db:"edit_state"`
		Amount          uint32        `db:"amount"`
		PayedAmount     uint32        `db:"payed_amount"`
		RefundAmount    uint32        `db:"refund_amount"`
		CreatedAt       time.Time     `db:"created_at"`
		UpdatedAt       time.Time     `db:"updated_at"`
		ReceiptItems    []ReceiptItem
		Payments        []Payment
	}

	Payment struct {
		Id            uint64        `db:"id"`
		OrderId       uint64        `db:"order_id"`
		Amount        uint32        `db:"amount"`
		PaymentMethod PaymentMethod `db:"payment_method"`
		PaymentLink   string        `db:"payment_link"`
		BankPaymentId string        `db:"bank_payment_id"`
		Payed         bool          `db:"payed"`
		Refunded      bool          `db:"refunded"`
		RefundAmount  uint32        `db:"refund_amount"`
		Expired       bool          `db:"expired"`
		CreatedAt     time.Time     `db:"created_at"`
		UpdatedAt     time.Time     `db:"updated_at"`
		PayedAt       sql.NullTime  `db:"payed_at"`
	}

	OrderMessage struct {
		Id          uint64      `db:"id"`
		OrderId     uint64      `db:"order_id"`
		DisplayMode DisplayMode `db:"display_mode"`
	}

	User struct {
		Id                   uint64    `db:"id"`
		IsMerchant           bool      `db:"is_merchant"`
		ActiveOrderId        uint64    `db:"active_order_id"`
		MerchantId           string    `db:"merchant_id"`
		SecretKey            string    `db:"secret_key"`
		ActiveOrderMessageId string    `db:"active_order_msg_id"`
		CreatedAt            time.Time `db:"created_at"`
		UpdatedAt            time.Time `db:"updated_at"`
	}
)

func (o Order) MessageString(c *Customer, mode DisplayMode) string {
	switch mode {
	case DisplayModeFull:
		return o.getFullMessageString(c)
	case DisplayModeCollapsed:
		return o.getCollapsedMessageString()
	case DisplayModeDeleted:
		return o.getDeletedMessageString()
	}

	return ""
}

func (o Order) getCollapsedMessageString() string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Now().Location()
	}
	createdAt := o.CreatedAt.In(loc).Format("01.02.2006")

	var payed string
	if o.Amount <= o.PayedAmount {
		payed = "Оплачен"
	} else {
		payed = fmt.Sprintf("Оплачено %d₽ из %d₽", o.PayedAmount/100, o.Amount/100)
	}

	orderState := o.OrderState.MessageString()
	result := fmt.Sprintf("*Заказ: #%d* от %s. %s.", o.Id, createdAt, orderState)
	if o.Amount > 0 {
		result = fmt.Sprintf("%s %s.", result, payed)
	}

	return result
}

func (o Order) getDeletedMessageString() string {
	return fmt.Sprintf("*Заказ #%d.* _Удалён_.", o.Id)
}

func (o Order) getFullMessageString(c *Customer) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Now().Location()
	}

	createdAt := o.CreatedAt.In(loc).Format("02.01.2006")

	amount := float32(o.Amount) / 100.0

	result := fmt.Sprintf(
		`*Заказ #%d* _(%s)_
*Создан:* %s
*Сумма:* %.2f₽
`, o.Id, o.OrderState.MessageString(), createdAt, amount)

	if o.Amount != 0 && o.PayedAmount != 0 {
		if o.PayedAmount >= o.Amount {
			result += "*Оплачен:* полностью\n"
		} else {
			payedAmount := float32(o.PayedAmount) / 100.0
			result += fmt.Sprintf("*Оплачено:* %.2f₽\n", payedAmount)
			restAmount := float32(o.Amount-o.PayedAmount) / 100.0
			result += fmt.Sprintf("*Остаток:* %.2f₽\n", restAmount)
		}
	}

	if o.RefundAmount != 0 {
		if o.RefundAmount >= o.Amount {
			result += "*Возврат:* в полном объеме\n"
		} else {
			refundAmount := float32(o.RefundAmount) / 100.0
			result += fmt.Sprintf("*Возвращено:* %.2f₽\n", refundAmount)
		}
	}

	result += "\n*Позиции*\n"

	if o.ReceiptItems != nil {
		for _, i := range o.ReceiptItems {
			price := float32(i.Price) / 100.0
			result += fmt.Sprintf("`- %s %.2f₽ x%d`\n", i.Name, price, i.Quantity)
		}
	}

	result += "\n*Клиент*"
	if c != nil {

		if c.Name.Valid {
			result += fmt.Sprintf("\n*Имя:* %s", c.Name.String)
		}

		if c.Email.Valid {
			email := strings.Replace(c.Email.String, "_", "\\_", -1)
			result += fmt.Sprintf("\n*E-Mail:* %s", email)
		} else {
			result += fmt.Sprintf("\n*E-Mail:* ‼️ Данные не заполнены")
		}

		if c.Phone.Valid {
			result += fmt.Sprintf("\n*Телефон:* %s", c.Phone.String)
		}

		if c.Instagram.Valid {
			result += fmt.Sprintf("\n*Instagram:* [@%[1]s](https://instagram.com/%[1]s)", c.Instagram.String)
		}

	} else {
		result += "\n‼️ Данные не заполнены"
	}

	result += "\n"

	if o.Payments != nil {
		result += "\n*Данные по оплате*"
		if len(o.Payments) == 0 {
			result += "\nНе найдено"
		} else {
			sort.Sort(PaymentsByCreatedAt(o.Payments))
			for i, p := range o.Payments {
				result += fmt.Sprintf("\n%s\n", p.MessageString(i+1))
			}
		}
	}

	return result
}

func (p Payment) MessageString(id int) string {
	result := fmt.Sprintf(`*Платеж #%d.* %s.`, id, p.PaymentMethod.MessageString(p.PaymentLink))

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Now().Location()
	}

	amount := float32(p.Amount) / 100.0
	result += fmt.Sprintf("\n*Сумма:* %.2f₽", amount)

	if p.PaymentMethod == PaymentMethodAcquiring {
		createdAt := p.CreatedAt.In(loc).Format("02 Jan 2006 15:04")
		result += fmt.Sprintf("\n*Создан:* %s", createdAt)

		if p.Payed && p.PayedAt.Valid {
			payedAt := p.PayedAt.Time.In(loc).Format("02 Jan 2006 15:04")
			result += fmt.Sprintf("\n*Оплачен:* %s", payedAt)
		} else {
			result += fmt.Sprintf("\n*Оплачен:* нет")
		}
	} else {
		payedAt := p.PayedAt.Time.In(loc).Format("02 Jan 2006 15:04")
		result += fmt.Sprintf("\n*Оплачен:* %s", payedAt)
	}

	if p.RefundAmount != 0 {
		refundAmount := float32(p.RefundAmount) / 100.0
		if p.RefundAmount == p.Amount {
			result += fmt.Sprintf("\n*Возвращен:* в полном объеме")
		} else {
			result += fmt.Sprintf("\n*Возвращено:* %.2f₽", refundAmount)
		}
	}

	return result
}

func (p PaymentMethod) MessageString(link string) string {
	switch p {
	case PaymentMethodCard2Card:
		return "Перевод на карту"
	case PaymentMethodAcquiring:
		return fmt.Sprintf("Оплата по [ссылке](%s)", link)
	case PaymentMethodCash:
		return "Оплата наличными"
	default:
		return "Неизвестный тип оплаты"
	}
}

func (s OrderState) MessageString() string {
	switch s {
	case OrderStateForming:
		return "Формируется"
	case OrderStateDone:
		return "Завершен"
	case OrderStateDeleted:
		return "Удалён"
	case OrderStateInProgress:
		return "В работе"
	}
	return "Статус неизвестен"
}

type PaymentsByCreatedAt []Payment

func (p PaymentsByCreatedAt) Len() int {
	return len(p)
}

func (p PaymentsByCreatedAt) Less(i, j int) bool {
	return p[i].CreatedAt.Before(p[j].CreatedAt)
}

func (p PaymentsByCreatedAt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
