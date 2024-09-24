package payment

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
)

// CreatePayment initiates a payment process by sending a request to the payment service.
// It sends an HTTP POST request with the payment details.
func CreatePayment(req interface{}) (*model.PaymentResponse, error) {
	var (
		paymentChannel = make(chan helper.Response)                                     // Channel to receive the response asynchronously
		paymentUri     = "https://29ac-110-136-183-86.ngrok-free.app/payments/payments" // URL of the payment service
	)

	// Prepare the HTTP client and request payload.
	payload := helper.NewNetClientRequest(paymentUri, helper.DefaultNetClient)
	payload.Post(req, paymentChannel) // Send the POST request with the payment data

	paymentRes := <-paymentChannel // Wait for the response from the channel

	// Check if there was an error in the response.
	if paymentRes.Err != nil {
		return nil, paymentRes.Err // Return the error if present
	}

	// Unmarshal the JSON response into the PaymentResponse struct.
	var response *model.PaymentResponse
	if err := json.Unmarshal(paymentRes.Res, &response); err != nil {
		return nil, err // Return an error if unmarshalling fails
	}

	return response, nil // Return the payment response on success
}
