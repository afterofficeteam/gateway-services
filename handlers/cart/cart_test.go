package cart

import (
	"bytes"
	"encoding/json"
	model "gateway-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupRequestAndResponse(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buf)
	rr := httptest.NewRecorder()
	return req, rr
}

func TestGetByUserID_Success(t *testing.T) {
	mockService := new(MockCartService)
	userID := uuid.New().String()
	reqBody := model.GetCartRequest{
		UserID:    uuid.MustParse(userID),
		ProductID: []uuid.UUID{uuid.New()},
	}

	expectedCart := []model.Cart{
		{
			ID:        uuid.New(),
			UserID:    uuid.MustParse(userID),
			ProductID: uuid.New(),
			Qty:       1,
		},
	}

	mockService.On("CartByUserID", reqBody, userID).Return(&expectedCart, nil)

	req, rr := setupRequestAndResponse(http.MethodGet, "/cart/user/"+userID, reqBody)

	handler := http.HandlerFunc(GetByUserID)
	handler.ServeHTTP(rr, req)

	var response model.Response
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
}
