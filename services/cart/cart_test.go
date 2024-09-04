package cart

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindByUserID(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	mockCart := []model.Cart{
		{
			ID:        uuid.New(),
			UserID:    userID,
			ProductID: productID,
			Qty:       2,
			CreatedAt: timePtr(time.Now()),
			UpdatedAt: timePtr(time.Now()),
		},
	}

	mockResponse, _ := json.Marshal(mockCart)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer ts.Close()

	req := model.GetCartRequest{
		UserID:    userID,
		ProductID: []uuid.UUID{productID},
	}

	helper.NetClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// Test FindByUserID function
	response, err := FindByUserID(req, userID.String())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
