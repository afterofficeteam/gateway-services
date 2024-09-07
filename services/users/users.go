package users

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
)

func Login(req model.UsersLogin) (*model.LoginResponse, error) {
	var (
		usersChannel = make(chan helper.Response)
		usersUri     = "url"
	)

	client := helper.NetClientRequest{
		NetClient:  helper.NetClient,
		RequestUrl: usersUri,
	}

	client.Post(req, usersChannel)

	usersRes := <-usersChannel
	if usersRes.Err != nil {
		return nil, usersRes.Err
	}

	var response *model.LoginResponse
	if err := json.Unmarshal(usersRes.Res, &response); err != nil {
		return nil, err
	}

	return response, nil
}
