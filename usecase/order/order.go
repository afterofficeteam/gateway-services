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

func UpdateOrder(req interface{}) (*string, error) {
	updateOK, err := order.UpdateOrder(req)
	if err != nil {
		return nil, err
	}

	return updateOK, nil
}

// CreateOrder handles the process of creating an order, validating product stock,
// updating stock levels after an order, and initiating payment via a bank transfer.
func CreateOrder(req model.PayloadCreateRequest) (*model.PaymentResponse, error) {
	// Step 1: Prepare product IDs for query
	// Extract product IDs from the request payload to fetch product details later.
	var productArr []string
	for _, p := range req.ProductOrder {
		productArr = append(productArr, p.ProductID) // Collect product IDs into a slice
	}
	productID := strings.Join(productArr, ",") // Join product IDs into a comma-separated string

	// Step 2: Find product details
	// Fetch the details of all products based on their IDs.
	product, err := productServices.FindDetail(productID)
	if err != nil {
		return nil, err // Return an error if product details cannot be retrieved
	}

	// Step 3: Validate stock and calculate prices
	// Ensure that there is sufficient stock for each product, and calculate the total price.
	if err := validateAndCalculatePrices(&req, *product); err != nil {
		return nil, err // Return an error if stock is insufficient or any other validation fails
	}

	// Step 4: Create order and get order ID
	// Create an order using the information from the request and product details.
	orderID, err := order.CreateOrder(req)
	if err != nil {
		return nil, err // Return an error if the order cannot be created
	}

	// Step 5: Update product stock after order
	// Reduce the stock of each product based on the quantities ordered.
	if err := updateProductStock(req, *product); err != nil {
		return nil, err // Return an error if stock update fails
	}

	// Step 6: Create payment
	// Initiate a payment process for the order via bank transfer.
	paymentOK, err := createPayment(req.BankTransfer, int(req.TotalPrice), *orderID)
	if err != nil {
		return nil, err // Return an error if payment initiation fails
	}

	// Return the payment response on success.
	return paymentOK, nil
}

// validateAndCalculatePrices ensures the requested product quantities are available,
// and computes the subtotal for each product and the total price for the order.
func validateAndCalculatePrices(req *model.PayloadCreateRequest, product model.DataProduct) error {
	for _, p := range product.Data.Items { // Loop over each product in the product data
		for i, r := range req.ProductOrder { // Loop over each product in the order request
			if p.Id == r.ProductID { // Match product ID from request with product data
				if p.Stock < r.Qty { // Check if requested quantity exceeds available stock
					return fmt.Errorf("stock is not enough") // Return an error if stock is insufficient
				}
				// Set price and calculate subtotal for each ordered product
				req.ProductOrder[i].Price = p.Price
				req.ProductOrder[i].SubtotalPrice = p.Price * float64(r.Qty)
			}
		}
	}

	// Calculate the total price for the order by summing up the subtotal prices.
	var total float64
	for _, p := range req.ProductOrder {
		total += p.SubtotalPrice // Add each product's subtotal to the total price
	}
	req.TotalPrice = total // Set the total price in the request

	return nil // Return nil if validation and calculation are successful
}

// updateProductStock updates the stock of each product after the order is successfully created.
func updateProductStock(req model.PayloadCreateRequest, product model.DataProduct) error {
	for _, r := range req.ProductOrder { // Loop over each product in the order request
		for _, p := range product.Data.Items { // Loop over each product in the product data
			if p.Id == r.ProductID { // Match product ID from request with product data
				// Prepare the payload for updating product stock
				payloadUpdate := struct {
					ProductId string `json:"product_id"`
					Stock     int    `json:"stock"`
				}{
					Stock:     p.Stock - r.Qty, // Reduce stock by the ordered quantity
					ProductId: p.Id,
				}
				// Append the update payload to the request's UpdateQty slice
				req.UpdateQty = append(req.UpdateQty, payloadUpdate)
			}
		}
	}

	// Call the service to update the product stock in the database
	_, err := productServices.Update(req.UpdateQty)
	if err != nil {
		return err // Return an error if stock update fails
	}

	return nil // Return nil if stock update is successful
}

// createPayment initiates a payment for the order via the specified bank transfer.
func createPayment(bank model.BankTransfer, totalPrice int, orderID string) (*model.PaymentResponse, error) {
	// Prepare the payload for creating the payment
	payloadPayment := struct {
		BankTransfer       model.BankTransfer       `json:"bank_transfer"`
		BasicAuthHeader    string                   `json:"basic_auth_header"`
		PaymentType        string                   `json:"payment_type"`
		TransactionDetails model.TransactionDetails `json:"transaction_details"`
	}{
		BankTransfer:    bank, // Bank transfer details from the request
		BasicAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("SB-Mid-server-jz-9ZTjDo8yA-5kZCU6rgDNr"+":")),
		PaymentType:     "bank_transfer", // Specify that payment is via bank transfer
		TransactionDetails: model.TransactionDetails{
			OrderID:     orderID,    // The order ID to associate with the payment
			GrossAmount: totalPrice, // The total price of the order
		},
	}

	// Call the payment service to initiate the payment
	paymentOK, err := payment.CreatePayment(payloadPayment)
	if err != nil {
		return nil, err // Return an error if payment creation fails
	}

	// Return the payment response on success
	return paymentOK, nil
}
