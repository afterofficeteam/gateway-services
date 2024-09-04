package cart

import (
	model "gateway-service/models"

	"github.com/stretchr/testify/mock"
)

type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) CartByUserID(req model.GetCartRequest, userID string) (*[]model.Cart, error) {
	args := m.Called(req, userID)
	return args.Get(0).(*[]model.Cart), args.Error(1)
}
