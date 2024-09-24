package order

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"gateway-service/usecase/order"
	"gateway-service/util/middleware"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateOrder is an HTTP handler that processes the creation of an order for a given user.
// This function validates the user ID, applies rate limiting, decodes the incoming request,
// and invokes the order creation process. It also handles response generation for success and error cases.
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Step 1: Extract user ID from the request path parameters.
	userID := r.PathValue("user_id") // Retrieve user ID from the URL path.
	if userID == "" {
		// If user ID is not provided, return a 400 Bad Request response with a relevant message.
		helper.HandleResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	// Step 2: Apply rate limiting based on user ID to prevent abuse (too many requests).
	if limiter := middleware.GetLimiter(userID); !limiter.Allow() {
		// If the request exceeds the allowed rate limit, return a 429 Too Many Requests response.
		helper.HandleResponse(w, http.StatusTooManyRequests, "Too many request, please try again later", nil)
		return
	}

	// Step 3: Convert the user ID from string to UUID format.
	usrIDcvt := uuid.MustParse(userID) // Convert userID string to UUID.

	// Step 4: Decode the request body into the PayloadCreateRequest struct.
	var bReq model.PayloadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		// If the request body cannot be decoded (invalid JSON), return a 400 Bad Request response.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 5: Set the converted userID into the request payload.
	bReq.UserID = usrIDcvt

	// Step 6: Call the CreateOrder function from the order service to create the order.
	bRes, err := order.CreateOrder(bReq)
	if err != nil {
		// If there is an error during order creation, return a 400 Bad Request response with the error message.
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Step 7: On successful order creation, send a 200 OK response with the created order details.
	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var bReq model.MidtransPayload
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if strings.Contains(bReq.StatusMessage, "notification") {
		helper.HandleResponse(w, http.StatusNotFound, "not a notif path", nil)
		return
	}

	timeNow := time.Now()
	payload := struct {
		OrderId   string     `json:"order_id"`
		Status    string     `json:"status"`
		IsPaid    bool       `json:"is_paid"`
		UpdatedAt *time.Time `json:"updated_at"`
	}{
		OrderId:   bReq.OrderID,
		Status:    "Payment",
		IsPaid:    true,
		UpdatedAt: &timeNow,
	}

	bRes, err := order.UpdateOrder(payload)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
