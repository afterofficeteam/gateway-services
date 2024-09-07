package cart

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"gateway-service/usecase/cart"
	"net/http"
)

func GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	var bReq model.GetCartRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bRes, err := cart.CartByUserID(bReq, userID)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
