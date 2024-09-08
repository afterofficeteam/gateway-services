package order

import (
	"encoding/base64"
	"fmt"
	model "gateway-service/models"
	"gateway-service/services/order"
	"gateway-service/services/payment"
	productServices "gateway-service/services/product"
	"strings"
)

func CreateOrder(req model.PayloadCreateRequest) (*string, error) {
	// Step 1: Prepare product IDs for query
	var productArr []string
	for _, p := range req.ProductOrder {
		productArr = append(productArr, p.ProductID)
	}
	productID := strings.Join(productArr, ",")

	// Step 2: Find product details
	product, err := productServices.FindDetail(productID)
	if err != nil {
		return nil, err
	}

	// Step 3: Validate stock and calculate prices
	if err := validateAndCalculatePrices(&req, *product); err != nil {
		return nil, err
	}

	// Step 4: Create order and get order ID
	orderID, err := order.CreateOrder(req)
	if err != nil {
		return nil, err
	}

	// Step 5: Update product stock after order
	if err := updateProductStock(req, *product); err != nil {
		return nil, err
	}

	// Step 6: Create payment
	if err := createPayment(req.TotalPrice, *orderID); err != nil {
		return nil, err
	}

	createOK := "Create order success"
	return &createOK, nil
}

// Step 3: Validate stock and calculate prices
func validateAndCalculatePrices(req *model.PayloadCreateRequest, product model.DataProduct) error {
	for _, p := range product.Data.Items {
		for i, r := range req.ProductOrder {
			if p.Id == r.ProductID {
				if p.Stock < r.Qty {
					return fmt.Errorf("stock is not enough")
				}
				req.ProductOrder[i].Price = p.Price
				req.ProductOrder[i].SubtotalPrice = p.Price * float64(r.Qty)
			}
		}
	}

	var total float64
	for _, p := range req.ProductOrder {
		total += p.SubtotalPrice
	}
	req.TotalPrice = total

	return nil
}

// Step 5: Update product stock after order
func updateProductStock(req model.PayloadCreateRequest, product model.DataProduct) error {
	for _, r := range req.ProductOrder {
		for _, p := range product.Data.Items {
			if p.Id == r.ProductID {
				payloadUpdate := struct {
					ProductId string `json:"product_id"`
					Stock     int    `json:"stock"`
				}{
					Stock:     p.Stock - r.Qty,
					ProductId: p.Id,
				}
				req.UpdateQty = append(req.UpdateQty, payloadUpdate)
			}
		}
	}

	_, err := productServices.Update(req.UpdateQty)
	return err
}

// Step 6: Create payment
func createPayment(totalPrice float64, orderID string) error {
	payloadPayment := struct {
		BasicAuthHeader string  `json:"basic_auth_header"`
		OrderID         string  `json:"order_id"`
		PaymentType     string  `json:"payment_type"`
		GrossAmount     float64 `json:"gross_amount"`
	}{
		BasicAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("SB-Mid-server-jz-9ZTjDo8yA-5kZCU6rgDNr"+":")),
		OrderID:         orderID,
		PaymentType:     "bank_transfer",
		GrossAmount:     totalPrice,
	}

	_, err := payment.CreatePayment(payloadPayment)
	return err
}
