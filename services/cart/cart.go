package cart

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
)

func FindByUserID(req model.GetCartRequest, userID string) (*[]model.Cart, error) {
	var (
		cartChannel = make(chan helper.Response)
		cartUri     = "http://localhost:9993/cart/" + userID
	)

	netClient := helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: cartUri,
	}
	netClient.Get(req, cartChannel)
	cartRes := <-cartChannel
	if cartRes.Err != nil {
		return nil, cartRes.Err
	}

	var response []model.Cart
	if err := json.Unmarshal(cartRes.Res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
