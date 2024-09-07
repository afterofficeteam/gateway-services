package users

import (
	model "gateway-service/models"
	"gateway-service/services/users"
)

func Login(req model.UsersLogin) (*model.LoginResponse, error) {
	return users.Login(req)
}
