package payment

import (
	model "gateway-service/models"
	"gateway-service/services/payment"
)

func CreatePayment(req model.PaymentRequest) (*model.PaymentResponse, error) {
	paymentOK, err := payment.CreatePayment(req)
	if err != nil {
		return nil, err
	}

	return paymentOK, nil
}
