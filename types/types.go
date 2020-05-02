package types

const (
	GOODS ReceiptItemType = iota
	DELIVERY
)

const (
	FULL_PREPAYMENT PaymentType = iota
	PARTIAL_PREPAYMENT
)

type (
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

	Item struct {
		Name  string
		Price uint32
		Type  ReceiptItemType
	}

	Customer struct {
		Id uint64
	}
)
