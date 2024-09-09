curl --location 'localhost:9990/gateway-service/order/037b4709-ee9c-4735-b213-9b82b02855ba' \
--header 'Content-Type: application/json' \
--data '{
    "product_order": [
        {
            "product_id": "c633d782-291d-4c85-bfe1-dfccf313d8f5",
            "qty": 1
        },
        {
            "product_id": "2e176973-19a0-41a7-bd50-b019eecc0e70",
            "qty": 1
        }
    ],
    "payment_type_id": "a09c0d3d-0d94-409a-8118-4bcc70cfc353",
    "order_number": "1",
    "status": "pending",
    "is_paid": false,
    "bank_transfer": {
        "bank": "bca"
    }
}'

curl --location 'localhost:9990/gateway-service/payment' \
--header 'Content-Type: application/json' \
--data '{
    "payment_type": "bank_transfer",
    "transaction_details": {
        "order_id": "e095f4e1-4c37-452a-ab02-939fe1ade69a",
        "gross_amount": 10000
    },
    "bank_transfer": {
        "bank": "bca"
    }
}'