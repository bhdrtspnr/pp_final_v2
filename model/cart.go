package model

import (
	"final_project/logger"
	connector "final_project/mysql"
	"time"
)

//define cartstructure for carts table
type Cart struct {
	Id            int       `json:"id"`
	CustomerId    int       `json:"customer_id"`
	IsPurchased   bool      `json:"is_purchased"`   //is used to check if cart is purchased
	DatePurchased time.Time `json:"date_purchased"` //date when cart is purchased
	TotalPrice    float64   `json:"total_price"`
	Discount      float64   `json:"discount"`
}

//check if a cart is completed
func IsCartCompleted(id string) bool {
	logger.AppLogger.Info().Println("Function hit : IsCartPurchased")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT is_purchased FROM carts WHERE id = ?", id) //query cart id and return is_purchased value
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var is_purchased bool
	for selDB.Next() {
		err = selDB.Scan(&is_purchased) //scan is_purchased value from cart id query result to is_purchased variable
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	logger.AppLogger.Info().Printf("Cart id: %s , is purchased: %t \n", id, is_purchased)
	return is_purchased //return is_purchased value
}

func IsProductExistsInCart(cart_id string, product_id string) bool {
	logger.AppLogger.Info().Println("Function hit : IsProductExistsInCart")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM cart_items WHERE cart_id = ? AND product_id = ?", cart_id, product_id) //query cart id and product id and return cart_item id
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() { //if cart_item id is found return true
		return true
	}
	return false //if cart_item id is not found return false
}

func UpdateCartPrice(id string) {
	//since I made a clear design error and did not hold product prices in the cart items I had to query cart_items to get ids of products and then query products to get prices :(
	//had to rewrite whole function due to that design error
	logger.AppLogger.Info().Println("Function hit : UpdateCartPrice")
	db := connector.DbConn()
	defer db.Close()
	//get current price
	curPrice, err := db.Query("SELECT total_price FROM carts WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var currentTotal float64
	for curPrice.Next() {
		err = curPrice.Scan(&currentTotal)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}

	//get cart items
	selDB, err := db.Query("SELECT PRODUCT_ID FROM cart_items WHERE cart_id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var product_ids []string //create slice to store product ids
	for selDB.Next() {
		var product_id string
		err = selDB.Scan(&product_id)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		product_ids = append(product_ids, product_id) //append product ids to slice
	}
	//get product prices
	var total_price float64
	for _, product_id := range product_ids {
		selDB, err := db.Query("SELECT price FROM products WHERE id = ?", product_id)
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
		total_price += price //add product prices to total price
	}
	//update total price
	_, err = db.Exec("UPDATE carts SET total_price = ? WHERE id = ?", total_price, id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error updating database: %v \n", err)
		panic(err.Error())
	}
	logger.AppLogger.Info().Printf("Cart id: %s price updated!, OLD VAL: %f , NEW_VAL: %f \n", id, currentTotal, total_price)

}

func GetCart(id string) Cart {
	//get cart from carts table by id
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
	return cart //return cart
}

func GetCartOwnerId(cart Cart) int {
	//get customer id from carts table by id
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
	return customer_id //return customer id
}

func CreateCart(cart Cart) {
	//create cart in carts table
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
	//check if cart exists in carts table
	logger.AppLogger.Info().Println("Function hit : IsCartExists")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM carts WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() {
		return true //cart exists
	}
	return false //cart does not exist
}
