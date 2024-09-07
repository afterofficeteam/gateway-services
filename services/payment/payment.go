package payment

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
)

func CreatePayment(req interface{}) (*model.CreatePaymentResponse, error) {
	var (
		paymentChannel = make(chan helper.Response)
		paymentUri     = "url"
	)

	payload := helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: paymentUri,
	}

	payload.Post(req, paymentChannel)
	paymentRes := <-paymentChannel

	if paymentRes.Err != nil {
		return nil, paymentRes.Err
	}

	var response *model.CreatePaymentResponse
	if err := json.Unmarshal(paymentRes.Res, &response); err != nil {
		return nil, err
	}

	return response, nil
}
