package domain

import "context"

//go:generate mockgen -source=ports.go -destination=./mocks/ports_mocks.go -package=mocks

type PaymentGateway interface {
    Charge(ctx context.Context, amountCents int64, currency, source string) (txID string, err error)
}

type OrderRepo interface {
    Save(ctx context.Context, o Order) error
}

type Order struct {
    ID          string
    AmountCents int64
    Currency    string
    Status      string // "pending", "paid", "failed"
    PaymentTxID string
}
