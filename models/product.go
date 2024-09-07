package model

import "time"

type UpdateQtyRequest struct {
	ProductId string `json:"product_id"`
	Stock     int    `json:"stock"`
}

type ProductOrder struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Price         float64 `json:"price"`
	Qty           int     `json:"qty"`
	SubtotalPrice float64 `json:"subtotal_price"`
}

type ProductResponse struct {
	Items []Product `json:"items"`
	Meta  Meta      `json:"meta"`
}

type Product struct {
	Id         string    `json:"id" db:"id"`
	CategoryId string    `json:"category_id" db:"category_id"`
	ShopId     string    `json:"shop_id" db:"shop_id"`
	Name       string    `json:"name" db:"name"`
	ImageUrl   *string   `json:"image_url" db:"image_url"`
	Price      float64   `json:"price" db:"price"`
	Stock      int       `json:"stock" db:"stock"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Meta struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}

type DataProduct struct {
	Data    ProductResponse `json:"data"`
	Message string          `json:"message"`
}
