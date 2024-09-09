package model

import (
	"time"

	"github.com/google/uuid"
)

type PayloadCreateRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`

	// Product Get
	ProductIds string `json:"product_ids"`

	// Product Update
	UpdateQty []UpdateQtyRequest `json:"update_qty"`

	// User
	PaymentTypeID uuid.UUID      `json:"payment_type_id" validate:"required"`
	OrderNumber   string         `json:"order_number" validate:"required"`
	TotalPrice    float64        `json:"total_price" validate:"required"`
	ProductOrder  []ProductOrder `json:"product_order"`
	Status        string         `json:"status" validate:"required"`
	IsPaid        bool           `json:"is_paid"`
	RefCode       string         `json:"ref_code"`
	CreatedAt     *time.Time     `json:"created_at"`

	// Payment
	BasicAuthHeader    string             `json:"basic_auth_header"`
	PaymentType        string             `json:"payment_type"`
	TransactionDetails TransactionDetails `json:"transaction_details"`
	BankTransfer       BankTransfer       `json:"bank_transfer"`
}

type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

type BankTransfer struct {
	Bank string `json:"bank"`
}

type MidtransPayload struct {
	TransactionTime        string `json:"transaction_time"`
	TransactionStatus      string `json:"transaction_status"`
	TransactionID          string `json:"transaction_id"`
	StatusMessage          string `json:"status_message"`
	StatusCode             string `json:"status_code"`
	SignatureKey           string `json:"signature_key"`
	PaymentType            string `json:"payment_type"`
	OrderID                string `json:"order_id"`
	MerchantID             string `json:"merchant_id"`
	MaskedCard             string `json:"masked_card"`
	GrossAmount            string `json:"gross_amount"`
	FraudStatus            string `json:"fraud_status"`
	ECI                    string `json:"eci"`
	Currency               string `json:"currency"`
	ChannelResponseMessage string `json:"channel_response_message"`
	ChannelResponseCode    string `json:"channel_response_code"`
	CardType               string `json:"card_type"`
	Bank                   string `json:"bank"`
	ApprovalCode           string `json:"approval_code"`
}
