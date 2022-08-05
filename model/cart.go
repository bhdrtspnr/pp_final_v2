package model

import (
	"final_project/logger"
	connector "final_project/mysql"
	"time"
)

type Cart struct {
	Id            int       `json:"id"`
	CustomerId    int       `json:"customer_id"`
	IsPurchased   bool      `json:"is_purchased"`
	DatePurchased time.Time `json:"date_purchased"`
	TotalPrice    float64   `json:"total_price"`
	Discount      float64   `json:"discount"`
}

func IsProductExistsInCart(cart_id string, product_id string) bool {
	logger.AppLogger.Info().Println("Function hit : IsProductExistsInCart")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM cart_items WHERE cart_id = ? AND product_id = ?", cart_id, product_id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() {
		return true
	}
	return false
}

func UpdateCartPrice(id string) {
	logger.AppLogger.Info().Println("Function hit : UpdateCartPrice")
	db := connector.DbConn()
	defer db.Close()
	selDB, err := db.Query("SELECT total_price FROM carts WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var total_price float64
	for selDB.Next() {
		err = selDB.Scan(&total_price)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}

	getItems, err := db.Prepare("SELECT PRODUCT_ID FROM cart_items WHERE CART_ID = ?")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	selDB, err = getItems.Query(id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var product_id int
	for selDB.Next() {
		err = selDB.Scan(&product_id)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	getPrice, err := db.Prepare("SELECT sum(PRICE) FROM products WHERE ID = ?")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	selDB, err = getPrice.Query(product_id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var price float64
	for selDB.Next() {
		err = selDB.Scan(&price)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	total_price += price
	_, err = db.Exec("UPDATE carts SET total_price = ? WHERE id = ?", total_price, id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error updating database: %v \n", err)
		panic(err.Error())
	}
	logger.AppLogger.Info().Printf("Cart id: %s , total price: %f \n", id, total_price)
}

func GetCart(id string) Cart {
	logger.AppLogger.Info().Println("Function hit : GetCart")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM carts WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var cart Cart
	for selDB.Next() {
		err = selDB.Scan(&cart.Id, &cart.CustomerId, &cart.IsPurchased, &cart.DatePurchased, &cart.TotalPrice, &cart.Discount)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	logger.AppLogger.Info().Printf("Cart id: %d , customer id: %d , is purchased: %t , date purchased: %v , total price: %f , discount: %f \n", cart.Id, cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice, cart.Discount)
	return cart
}

func GetCartOwner(cart Cart) int {
	logger.AppLogger.Info().Println("Function hit : GetCartOwner")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT customer_id FROM carts WHERE id = ?", cart.Id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var customer_id int
	for selDB.Next() {
		err = selDB.Scan(&customer_id)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	logger.AppLogger.Info().Printf("Cart id: %d , customer id: %d \n", cart.Id, customer_id)
	return customer_id
}

func CreateCart(cart Cart) {
	logger.AppLogger.Info().Println("Function hit : CreateCart")
	db := connector.DbConn()
	insForm, err := db.Prepare("INSERT INTO carts(customer_id, is_purchased, date_purchased, total_price, discount) VALUES(?,?,?,?,?)")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	insForm.Exec(cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice, cart.Discount)
	logger.AppLogger.Info().Printf("Cart created with id: %d , customer id: %d , is purchased: %t , date purchased: %v , total price: %f , discount: %f \n", cart.Id, cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice, cart.Discount)
	defer db.Close()
}

func IsCartExists(id string) bool {
	logger.AppLogger.Info().Println("Function hit : IsCartExists")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM carts WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() {
		return true
	}
	return false
}
