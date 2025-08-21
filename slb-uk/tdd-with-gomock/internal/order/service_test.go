package order_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/slb-uk/tdd-with-gomock/internal/domain"
	"github.com/slb-uk/tdd-with-gomock/internal/domain/mocks"
	"github.com/slb-uk/tdd-with-gomock/internal/order"
)

func TestService_PlaceOrder_Success(t *testing.T) {
    t.Parallel()

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockPay := mocks.NewMockPaymentGateway(ctrl)
    mockRepo := mocks.NewMockOrderRepo(ctrl)

    svc := order.NewService(mockPay, mockRepo)

    in := domain.Order{ID: "ord_1", AmountCents: 4999, Currency: "INR", Status: "pending"}
    source := "tok_visa"

    mockPay.EXPECT().
        Charge(gomock.Any(), int64(4999), "INR", source).
        Return("tx_abc123", nil).
        Times(1)

    mockRepo.EXPECT().
        Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Order{})).
        Return(nil).
        Times(1)

    out, err := svc.PlaceOrder(context.Background(), in, source)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out.Status != "paid" || out.PaymentTxID != "tx_abc123" {
        t.Fatalf("unexpected order: %+v", out)
    }
}

func TestService_PlaceOrder_Variants(t *testing.T) {
    t.Parallel()

    cases := []struct {
        name          string
        chargeErr     error
        saveErr       error
        wantStatus    string
        wantErrSubstr string
    }{
        {"success", nil, nil, "paid", ""},
        {"charge fails", errors.New("card declined"), nil, "failed", "charge failed"},
        {"save fails", nil, errors.New("db down"), "paid", "save failed"},
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockPay := mocks.NewMockPaymentGateway(ctrl)
            mockRepo := mocks.NewMockOrderRepo(ctrl)
            svc := order.NewService(mockPay, mockRepo)

            in := domain.Order{ID: "ord_1", AmountCents: 4999, Currency: "INR", Status: "pending"}
            source := "tok_visa"

            mockPay.EXPECT().
                Charge(gomock.Any(), int64(4999), "INR", source).
                Return("tx_ok", tc.chargeErr).
                Times(1)

            if tc.chargeErr == nil {
                mockRepo.EXPECT().
                    Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Order{})).
                    Return(tc.saveErr).
                    Times(1)
            }

            out, err := svc.PlaceOrder(context.Background(), in, source)

            if tc.wantErrSubstr == "" && err != nil {
                t.Fatalf("unexpected err: %v", err)
            }
            if tc.wantErrSubstr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErrSubstr)) {
                t.Fatalf("want err containing %q, got %v", tc.wantErrSubstr, err)
            }
            if out.Status != tc.wantStatus {
                t.Fatalf("want status %q, got %q", tc.wantStatus, out.Status)
            }
        })
    }
}
