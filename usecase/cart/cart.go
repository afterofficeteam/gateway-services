package cart

import (
	model "gateway-service/models"
	"gateway-service/services/cart"
)

func CartByUserID(req model.GetCartRequest, userID string) (*[]model.Cart, error) {
	return cart.FindByUserID(req, userID)
}
