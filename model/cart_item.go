package model

type CartItem struct {
	Id          int    `json:"id"`
	CartId      int    `json:"cart_id"`
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`
}
