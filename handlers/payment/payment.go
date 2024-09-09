package payment

import (
	"encoding/base64"
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"gateway-service/usecase/payment"
	"net/http"
)

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	var bReq model.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bReq.BasicAuthHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte("SB-Mid-server-jz-9ZTjDo8yA-5kZCU6rgDNr"+":"))

	bRes, err := payment.CreatePayment(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusCreated, helper.SUCCESS_MESSSAGE, bRes)
}
