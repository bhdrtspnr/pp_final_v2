package rest_functions

import (
	"encoding/json"
	logger "final_project/logger"
	model "final_project/model"
	connector "final_project/mysql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var STATUS_NOT_ALLOWED int = 405
var STATUS_NOT_FOUND int = 404

func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	db := connector.DbConn()
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(STATUS_NOT_ALLOWED), http.StatusMethodNotAllowed)
		w.Write([]byte("Method: DELETE is required."))
		logger.AppLogger.Error().Printf("Method Delete is required. \n")
		return
	}
	params := mux.Vars(r) //get the params from the url
	cartId := params["cart_id"]
	productId := params["product_id"]
	if !model.IsCartExists(cartId) {
		http.Error(w, http.StatusText(STATUS_NOT_FOUND), http.StatusNotFound)
		w.Write([]byte("Cart: " + cartId + " not found"))
		logger.AppLogger.Error().Printf("Cart: %s not found", cartId)
		return
	}
	if !model.IsProductExists(productId) {
		http.Error(w, http.StatusText(STATUS_NOT_FOUND), http.StatusNotFound)
		w.Write([]byte("Product: " + productId + " not found in the DB"))
		logger.AppLogger.Error().Printf("Product not found in DB: %s", productId)
		return
	}

	if !model.IsProductExistsInCart(cartId, productId) {
		http.Error(w, http.StatusText(STATUS_NOT_FOUND), http.StatusNotFound)
		w.Write([]byte("Product: " + productId + " not found in cart: " + cartId))
		logger.AppLogger.Error().Printf("Product: %s not found in cart: %s", productId, cartId)
		return
	}

	logger.AppLogger.Info().Printf("Deleting product id: %s from cart id: %s \n", productId, cartId)
	delForm, err := db.Prepare("DELETE FROM cart_items WHERE cart_id=? AND product_id=? LIMIT 1")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	delForm.Exec(cartId, productId)
	logger.AppLogger.Info().Println("Product deleted from cart.")
	model.UpdateCartPrice(cartId)
	w.Write([]byte("Product: " + productId + " deleted from cart: " + cartId))
	defer db.Close()
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	logger.AppLogger.Info().Println("Endpoint hit : ShowCart")
	db := connector.DbConn()
	params := mux.Vars(r) //get the params from the url
	cartId := params["cart_id"]
	if !model.IsCartExists(cartId) {
		logger.AppLogger.Error().Println("Cart not found")
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	selCart, err := db.Query("SELECT * FROM carts WHERE id = ?", cartId)
	if err != nil {
		logger.AppLogger.Error().Println(err.Error())
		panic(err.Error())
	}
	cart := model.Cart{}
	for selCart.Next() {
		var id int
		var customerId int
		var isPurchased bool
		var DatePurchased time.Time
		var TotalPrice float64
		var Discount float64

		err = selCart.Scan(&id, &customerId, &isPurchased, &DatePurchased, &TotalPrice, &Discount)
		if err != nil {
			logger.AppLogger.Error().Println(err.Error())
			panic(err.Error())
		}
		cart.Id = id
		cart.CustomerId = customerId
		cart.IsPurchased = isPurchased
		cart.DatePurchased = DatePurchased
		cart.TotalPrice = TotalPrice
		cart.Discount = Discount
		logger.AppLogger.Info().Printf("Reading cart id: %d , customer id: %d , is purchased: %v , date purchased: %v , total price: %f \n", cart.Id, cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice)
	}
	json.NewEncoder(w).Encode(cart)

	selDB, err := db.Query("SELECT * FROM cart_items WHERE cart_id = ?", cartId)
	if err != nil {
		logger.AppLogger.Error().Println(err.Error())
	}
	cartItem := model.CartItem{}
	cartItems := []model.CartItem{}
	logger.AppLogger.Info().Println("Parsing cart items...")
	for selDB.Next() {
		var id int
		var cartId int
		var productId int
		var productName string

		err = selDB.Scan(&id, &cartId, &productId, &productName)
		if err != nil {
			logger.AppLogger.Error().Println(err.Error())
			panic(err.Error())
		}
		cartItem.Id = id
		cartItem.CartId = cartId
		cartItem.ProductId = productId
		cartItem.ProductName = productName
		cartItems = append(cartItems, cartItem)
		logger.AppLogger.Info().Printf("Reading cart item id: %d , cart id: %d , product id: %d , product name: %s \n", cartItem.Id, cartItem.CartId, cartItem.ProductId, cartItem.ProductName)
	}
	json.NewEncoder(w).Encode(cartItems)

	logger.AppLogger.Info().Printf("Total of %v cart items found and listed. \n", len(cartItems))
	defer db.Close()
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	db := connector.DbConn()
	if r.Method != "POST" {
		logger.AppLogger.Error().Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r) //get the params from the url
	cartId := params["cart_id"]
	productId := params["product_id"]
	productName := model.GetProductName(productId)

	if !model.IsCartExists(cartId) {
		logger.AppLogger.Error().Println("Cart not found")
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	if !model.IsProductExists(productId) {
		logger.AppLogger.Error().Println("Product not found")
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	logger.AppLogger.Info().Printf("Adding product id: %s , product name: %s to cart id: %s \n", productId, productName, cartId)
	insForm, err := db.Prepare("INSERT INTO cart_items (cart_id, product_id, product_name) VALUES(?,?,?)")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	insForm.Exec(cartId, productId, productName)
	logger.AppLogger.Info().Println("Product added to cart.")
	model.UpdateCartPrice(cartId)
	defer db.Close()
}

func ListCustomers(w http.ResponseWriter, r *http.Request) {
	logger.AppLogger.Info().Println("Endpoint hit : ListCustomers")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM customers")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	customer := model.Customer{}
	customers := []model.Customer{}
	logger.AppLogger.Info().Println("Parsing customers...")
	for selDB.Next() {
		var id int
		var name string
		var balance float64
		var consecutive_discount int
		var has_subsequent_discount_until time.Time
		err = selDB.Scan(&id, &name, &balance, &consecutive_discount, &has_subsequent_discount_until)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		customer.Id = id
		customer.Name = name
		customer.Balance = balance
		customer.Consecutive_discount = consecutive_discount
		customer.Has_subsequent_discount_until = has_subsequent_discount_until
		logger.AppLogger.Info().Printf("Reading customer id: %d , name: %s , balance: %f , consecutive discount: %d , has subsequent discount until: %v \n", customer.Id, customer.Name, customer.Balance, customer.Consecutive_discount, customer.Has_subsequent_discount_until)
		customers = append(customers, customer)
	}
	json.NewEncoder(w).Encode(customers)
	logger.AppLogger.Info().Printf("Total of %v customers found and listed. \n", len(customers))
	defer db.Close()
}

func ListCartItems(w http.ResponseWriter, r *http.Request) {
	logger.AppLogger.Info().Println("Endpoint hit : ListCartItems")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM cart_items")
	if err != nil {
		logger.AppLogger.Error().Println(err.Error())
	}
	cartItem := model.CartItem{}
	cartItems := []model.CartItem{}
	logger.AppLogger.Info().Println("Parsing cart items...")
	for selDB.Next() {
		var id int
		var cartId int
		var productId int
		var productName string

		err = selDB.Scan(&id, &cartId, &productId, &productName)
		if err != nil {
			logger.AppLogger.Error().Println(err.Error())
		}
		cartItem.Id = id
		cartItem.CartId = cartId
		cartItem.ProductId = productId
		cartItem.ProductName = productName
		cartItems = append(cartItems, cartItem)
		logger.AppLogger.Info().Printf("Reading cart item id: %d , cart id: %d , product id: %d , product name: %s \n", cartItem.Id, cartItem.CartId, cartItem.ProductId, cartItem.ProductName)
	}
	json.NewEncoder(w).Encode(cartItems)
	logger.AppLogger.Info().Printf("Total of %v cart items found and listed. \n", len(cartItems))
	defer db.Close()
}

func ListCarts(w http.ResponseWriter, r *http.Request) {
	logger.AppLogger.Info().Println("Endpoint hit : ListCarts")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM carts")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	cart := model.Cart{}
	carts := []model.Cart{}
	logger.AppLogger.Info().Println("Parsing carts...")
	for selDB.Next() {
		var id int
		var customer_id int
		var is_purchased bool
		var date_purchased time.Time
		var total_price float64
		var discount float64

		err = selDB.Scan(&id, &customer_id, &is_purchased, &date_purchased, &total_price, &discount)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		cart.Id = id
		cart.CustomerId = customer_id
		cart.IsPurchased = is_purchased
		cart.DatePurchased = date_purchased
		cart.TotalPrice = total_price
		carts = append(carts, cart)
		logger.AppLogger.Info().Printf("Reading cart id: %d , customer id: %d , is purchased: %t , date purchased: %v , total price: %f , discount: %f \n", cart.Id, cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice, cart.Discount)
	}
	logger.AppLogger.Info().Println("Marshalling carts...")
	json.NewEncoder(w).Encode(carts)
	logger.AppLogger.Info().Printf("Total of %v carts found and listed. \n", len(carts))
	defer db.Close()

}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	logger.AppLogger.Info().Println("Endpoint hit : ListProducts")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM products")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	product := model.Product{}
	products := []model.Product{}
	logger.AppLogger.Info().Println("Parsing products...")
	for selDB.Next() {
		var id int
		var name string
		var price float64
		var vat float64

		err = selDB.Scan(&id, &name, &price, &vat)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		logger.AppLogger.Info().Printf("Reading Product, ID: %v , name: %v , price: %v , vat: %v \n", id, name, price, vat)
		product.Id = id
		product.Name = name
		product.Price = price
		product.Vat = vat
		products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
	logger.AppLogger.Info().Printf("Total of %v products found and listed. \n", len(products))
	defer db.Close()
}
