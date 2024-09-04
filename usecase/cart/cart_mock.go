package cart

import (
	model "gateway-service/models"

	"github.com/stretchr/testify/mock"
)

type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) FindByUserID(req model.GetCartRequest, userID string) (*[]model.Cart, error) {
	args := m.Called(req, userID)
	if carts, ok := args.Get(0).(*[]model.Cart); ok {
		return carts, args.Error(1)
	}
	return nil, args.Error(1)
}
