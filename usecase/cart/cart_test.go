package cart

import (
	model "gateway-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCartByUserID(t *testing.T) {
	mockService := new(MockCartService)

	userID := uuid.New().String()
	req := model.GetCartRequest{
		UserID:    uuid.MustParse(userID),
		ProductID: []uuid.UUID{uuid.New()},
	}

	expectedCarts := &[]model.Cart{
		{
			ID:        uuid.New(),
			UserID:    uuid.MustParse(userID),
			ProductID: uuid.New(),
			Qty:       2,
			CreatedAt: nil,
			UpdatedAt: nil,
			DeletedAt: nil,
		},
	}

	mockService.On("FindByUserID", req, userID).Return(expectedCarts, nil)

	carts, err := CartByUserID(req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, carts)
}
