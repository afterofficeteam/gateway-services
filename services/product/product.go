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
		productUri     = "https://f45f-36-72-214-46.ngrok-free.app/products/products"
	)

	headerPayload := helper.NetClientRequest{
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

	headerPayload.Get(nil, productChannel)
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
		updateUrl     = "https://f45f-36-72-214-46.ngrok-free.app/products/product-stocks"
	)

	mutex.Lock()
	defer mutex.Unlock()

	header := helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: updateUrl,
	}

	header.Patch(req, updateChannel)
	updateRes := <-updateChannel

	if updateRes.Err != nil {
		return nil, updateRes.Err
	}

	response := string(updateRes.Res)

	return &response, nil
}
