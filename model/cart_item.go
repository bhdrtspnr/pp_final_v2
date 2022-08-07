package model

import (
	"final_project/logger"
	connector "final_project/mysql"
)

type CartItem struct {
	Id          int    `json:"id"`
	CartId      int    `json:"cart_id"`
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`
}

func GetCartItems(cartid string) []CartItem {
	logger.AppLogger.Info().Println("Function hit : GetCartItems")
	db := connector.DbConn()
	defer db.Close()
	selDB, err := db.Query("SELECT * FROM cart_items WHERE cart_id = ?", cartid)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var cartitems []CartItem
	for selDB.Next() {
		var cartitem CartItem
		err = selDB.Scan(&cartitem.Id, &cartitem.CartId, &cartitem.ProductId, &cartitem.ProductName)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		cartitems = append(cartitems, cartitem)
	}
	return cartitems
}
