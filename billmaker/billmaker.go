package billmaker

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/types"
)

type (
	BillMaker interface {
		CreateBill(ctx context.Context, order types.Order) (types.Payment, error)
	}

	billMaker struct {
	}
)

func NewBillMaker() BillMaker {
	return billMaker{}
}

func (bm billMaker) CreateBill(ctx context.Context, order types.Order) (types.Payment, error) {
	bill := types.Payment{
		Id: 1337,
		//Url: "https://modulbank.ru",
	}
	return bill, nil
}
