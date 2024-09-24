package order

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"sync"
)

var (
	mutex sync.Mutex
)

// CreateOrder sends a request to the order service to create a new order.
// It sends an HTTP POST request with the order details.
func CreateOrder(req model.PayloadCreateRequest) (*string, error) {
	var (
		orderChannel = make(chan helper.Response)                              // Channel to receive the response asynchronously
		orderUri     = "http://localhost:9993/cart-order-service/order/create" // URL of the order service
	)

	mutex.Lock()         // Lock the mutex to prevent concurrent writes
	defer mutex.Unlock() // Ensure the mutex is unlocked after the function completes

	// Prepare the HTTP client and request payload.
	payload := helper.NewNetClientRequest(orderUri, helper.DefaultNetClient)
	payload.Post(req, orderChannel) // Send the POST request with the order data

	orderRes := <-orderChannel // Wait for the response from the channel

	// Check if there was an error in the response.
	if orderRes.Err != nil {
		return nil, orderRes.Err // Return the error if present
	}

	// Unmarshal the JSON response into a string (order ID).
	var response *string
	if err := json.Unmarshal(orderRes.Res, &response); err != nil {
		return nil, err // Return an error if unmarshalling fails
	}

	return response, nil // Return the order ID on success
}

func UpdateOrder(req interface{}) (*string, error) {
	var (
		orderChannel   = make(chan helper.Response)
		updateOrderUri = "http://localhost:9993/cart-order-service/order/callback"
	)

	mutex.Lock()
	defer mutex.Unlock()

	payload := helper.NewNetClientRequest(updateOrderUri, helper.DefaultNetClient)
	payload.Post(req, orderChannel)

	orderRes := <-orderChannel
	if orderRes.Err != nil {
		return nil, orderRes.Err
	}

	var response string
	if err := json.Unmarshal(orderRes.Res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
