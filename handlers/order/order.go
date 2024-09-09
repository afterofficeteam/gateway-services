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

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	if limiter := middleware.GetLimiter(userID); !limiter.Allow() {
		helper.HandleResponse(w, http.StatusTooManyRequests, "Too many request, please try again request", nil)
		return
	}

	usrIDcvt := uuid.MustParse(userID)

	var bReq model.PayloadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bReq.UserID = usrIDcvt

	bRes, err := order.CreateOrder(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

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
