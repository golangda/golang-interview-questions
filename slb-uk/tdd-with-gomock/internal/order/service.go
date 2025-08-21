package order

import (
	"context"
	"fmt"

	"github.com/slb-uk/tdd-with-gomock/internal/domain"
)

type Service struct {
    pay domain.PaymentGateway
    db  domain.OrderRepo
}

func NewService(pay domain.PaymentGateway, db domain.OrderRepo) *Service {
    return &Service{pay: pay, db: db}
}

// PlaceOrder charges and, if successful, persists the order.
func (s *Service) PlaceOrder(ctx context.Context, o domain.Order, source string) (domain.Order, error) {
    if o.AmountCents <= 0 {
        return o, fmt.Errorf("invalid amount")
    }

    txID, err := s.pay.Charge(ctx, o.AmountCents, o.Currency, source)
    if err != nil {
        o.Status = "failed"
        return o, fmt.Errorf("charge failed: %w", err)
    }

    o.Status = "paid"
    o.PaymentTxID = txID

    if err := s.db.Save(ctx, o); err != nil {
        return o, fmt.Errorf("save failed: %w", err)
    }

    return o, nil
}
