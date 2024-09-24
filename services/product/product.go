package product

import (
	"encoding/json"
	"gateway-service/helper"
	model "gateway-service/models"
	"sync"
)

var (
	mutex sync.Mutex
)

// FindDetail retrieves the product details for the given product IDs.
// It sends an HTTP GET request to the product service and parses the response.
func FindDetail(productID string) (res *model.DataProduct, err error) {
	var (
		productChannel = make(chan helper.Response)                                     // Channel to receive the response asynchronously
		productUri     = "https://8a43-110-136-183-86.ngrok-free.app/products/products" // URL of the product service
	)

	// Prepare the request payload with necessary query parameters.
	headerPayload := helper.NetClientRequest{
		NetClient:  helper.DefaultNetClient, // Use the default HTTP client
		RequestUrl: productUri,              // Set the request URL
		QueryParam: []helper.QueryParams{ // Add query parameters to the request
			{
				Param: "product_ids",
				Value: productID, // Product IDs to fetch details for
			},
			{
				Param: "limit",
				Value: "100", // Limit the number of products returned
			},
		},
	}

	// Send the GET request to the product service.
	headerPayload.Get(nil, productChannel) // Send the request without additional headers
	productRes := <-productChannel         // Wait for the response from the channel

	// Check if there was an error in the response.
	if productRes.Err != nil {
		return nil, productRes.Err // Return the error if present
	}

	// Unmarshal the JSON response into the DataProduct struct.
	var response *model.DataProduct
	if err := json.Unmarshal(productRes.Res, &response); err != nil {
		return nil, err // Return an error if unmarshalling fails
	}

	return response, nil // Return the product details on success
}

// Update sends a request to update the stock quantities of products.
// It sends an HTTP PATCH request with the updated stock information.
func Update(req []model.UpdateQtyRequest) (*string, error) {
	var (
		updateChannel = make(chan helper.Response)                                         // Channel to receive the response asynchronously
		updateUrl     = "https://f45f-36-72-214-46.ngrok-free.app/products/product-stocks" // URL for updating product stocks
	)

	mutex.Lock()         // Lock the mutex to prevent concurrent writes
	defer mutex.Unlock() // Ensure the mutex is unlocked after the function completes

	// Prepare the HTTP client and request payload.
	payload := helper.NewNetClientRequest(updateUrl, helper.DefaultNetClient)
	payload.Patch(req, updateChannel) // Send the PATCH request with the update data

	updateRes := <-updateChannel // Wait for the response from the channel

	// Check if there was an error in the response.
	if updateRes.Err != nil {
		return nil, updateRes.Err // Return the error if present
	}

	// Convert the response bytes to a string.
	updateOK := string(updateRes.Res)

	return &updateOK, nil // Return the update confirmation message
}
