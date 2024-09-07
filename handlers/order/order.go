package order

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"gateway-service/usecase/order"
	"net/http"

	"github.com/google/uuid"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	usrIDcvt := uuid.MustParse(userID)

	param := r.URL.Query()
	limit := param.Get("limit")
	if limit == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "Limit is required", nil)
		return
	}

	var bReq model.PayloadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bReq.UserID = usrIDcvt
	bReq.Limit = limit

	bRes, err := order.CreateOrder(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
