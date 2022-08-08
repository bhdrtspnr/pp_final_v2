package model

import (
	"final_project/logger"
	connector "final_project/mysql"
)

//define cart_item structure for cart_items table
type CartItem struct {
	Id          int    `json:"id"`
	CartId      int    `json:"cart_id"`
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`
}

//get all the cart items for a cart with given cart_id
func GetCartItems(cartid string) []CartItem {
	logger.AppLogger.Info().Println("Function hit : GetCartItems")
	db := connector.DbConn()
	defer db.Close()
	selDB, err := db.Query("SELECT * FROM cart_items WHERE cart_id = ?", cartid)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var cartitems []CartItem //create a slice of cart items
	for selDB.Next() {
		var cartitem CartItem
		err = selDB.Scan(&cartitem.Id, &cartitem.CartId, &cartitem.ProductId, &cartitem.ProductName)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		cartitems = append(cartitems, cartitem) //append the cart item to the slice
	}
	return cartitems //return the slice of cart items
}
