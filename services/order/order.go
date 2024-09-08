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

func CreateOrder(req model.PayloadCreateRequest) (*string, error) {
	var (
		orderChannel = make(chan helper.Response)
		orderUri     = "http://localhost:9993/cart-order-service/order/create"
	)

	mutex.Lock()
	defer mutex.Unlock()

	client := &helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: orderUri,
	}

	client.Post(req, orderChannel)

	orderRes := <-orderChannel
	if orderRes.Err != nil {
		return nil, orderRes.Err
	}

	var response *string
	if err := json.Unmarshal(orderRes.Res, &response); err != nil {
		return nil, err
	}

	return response, nil
}
