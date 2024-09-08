package product

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"sync"
)

var (
	mutex sync.Mutex
)

func FindDetail(productID string) (res *model.DataProduct, err error) {
	var (
		productChannel = make(chan helper.Response)
		productUri     = "https://2421-36-72-214-46.ngrok-free.app/products/products"
	)

	payload := helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: productUri,
		QueryParam: []helper.QueryParams{
			{
				Param: "product_ids",
				Value: productID,
			},
			{
				Param: "limit",
				Value: "100",
			},
		},
	}

	payload.Get(nil, productChannel)
	productRes := <-productChannel

	if productRes.Err != nil {
		return nil, productRes.Err
	}

	var response *model.DataProduct
	if err := json.Unmarshal(productRes.Res, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func Update(req []model.UpdateQtyRequest) (*string, error) {
	var (
		updateChannel = make(chan helper.Response)
		updateUrl     = "https://2421-36-72-214-46.ngrok-free.app/products/product-stocks"
	)

	mutex.Lock()
	defer mutex.Unlock()

	payload := helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: updateUrl,
	}

	payload.Patch(req, updateChannel)
	updateRes := <-updateChannel

	if updateRes.Err != nil {
		return nil, updateRes.Err
	}

	var response *string
	if err := json.Unmarshal(updateRes.Res, &response); err != nil {
		return nil, err
	}

	return response, nil
}
