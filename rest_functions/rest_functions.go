package rest_functions

import (
	"encoding/json"
	"final_project/config"
	logger "final_project/logger"
	model "final_project/model"
	connector "final_project/mysql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	discount "final_project/discount"

	"github.com/gorilla/mux"
)

//IDE kept giving warnings about using the same error code again and again so I used a variable to store the error codes.
var STATUS_NOT_ALLOWED int = 405
var STATUS_NOT_FOUND int = 404

func CompleteCart(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Endpoint hit: complete cart \n")
	logger.AppLogger.Info().Println("Endpoint Hit: CompleteCart")
	db := connector.DbConn() //connect to the database
	if r.Method != "POST" {  //check if the method is POST
		logger.AppLogger.Error().Printf("Method %s is not allowed at endpoint complete cart.", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)            //get the params from the url
	cartId := params["cart_id"]      //get the cart id from the url
	if !model.IsCartExists(cartId) { //check if the cart id exists in the database
		logger.AppLogger.Error().Printf("Cart: %s not found\n", cartId)
		http.Error(w, "Cart not found with id ", http.StatusNotFound)
		return
	}

	if model.IsCartCompleted(cartId) { //check if the cart is completed
		logger.AppLogger.Error().Printf("Cart: %s is already completed\n", cartId)
		http.Error(w, "Cart is already completed", http.StatusBadRequest)
		return
	}

	discountAmount := 0.0                                                                        //initialize the discount amount to 0
	consecutivePurchaseDiscount := discount.CalculateConsecutivePurchaseDiscount(cartId)         //calculate the consecutive purchase discount
	givenAmountDiscount := discount.CalculateGivenAmountDiscount(cartId)                         //calculate the given amount discount
	ThreeSubsequentPurchaseDiscount := discount.CalculateThreeSubsequentPurchaseDiscount(cartId) //calculate the three subsequent purchase discount

	logger.AppLogger.Info().Printf("Comparing the discounts: consecutive purchase discount: %f, given amount discount: %f, three subsequent purchase discount: %f \n", consecutivePurchaseDiscount, givenAmountDiscount, ThreeSubsequentPurchaseDiscount)

	/* compare the discounts and get the highest one */
	if consecutivePurchaseDiscount > givenAmountDiscount && consecutivePurchaseDiscount > ThreeSubsequentPurchaseDiscount {
		logger.AppLogger.Info().Printf("Using consecutive purchase discount of: %f\n", consecutivePurchaseDiscount)
		discountAmount = consecutivePurchaseDiscount
	} else if givenAmountDiscount > consecutivePurchaseDiscount && givenAmountDiscount > ThreeSubsequentPurchaseDiscount {
		logger.AppLogger.Info().Printf("Using given amount discount of: %f\n", givenAmountDiscount)
		discountAmount = givenAmountDiscount
	} else if ThreeSubsequentPurchaseDiscount > consecutivePurchaseDiscount && ThreeSubsequentPurchaseDiscount > givenAmountDiscount {
		logger.AppLogger.Info().Printf("Using three subsequent purchase discount of: %f\n", ThreeSubsequentPurchaseDiscount)
		discountAmount = ThreeSubsequentPurchaseDiscount
	}

	/* end of comparing the discounts */

	logger.AppLogger.Info().Printf("Discount: %f \n", discountAmount)
	totalPrice := model.GetCart(cartId).TotalPrice - discountAmount //get the total price of the cart and subtract the discount amount from it
	logger.AppLogger.Info().Printf("Cart: %s total price: %f \n", cartId, totalPrice)
	if totalPrice < 0 { //check if the total price is less than 0
		totalPrice = 0
	}
	logger.AppLogger.Info().Printf("Cart: %s total price: %f \n", cartId, totalPrice)

	//check if user has enough balance to complete the cart
	if model.GetCustomer(model.GetCartOwnerId(model.GetCart(cartId))).Balance < totalPrice {
		logger.AppLogger.Error().Printf("Customer: %d does not have enough balance to complete the cart: %s \n", model.GetCartOwnerId(model.GetCart(cartId)), cartId)
		http.Error(w, "Customer does not have enough balance to complete the cart", http.StatusNotFound)
		return
	}

	//update the cart status to completed
	completeForm, err := db.Prepare("UPDATE carts SET is_purchased=1, total_price=? WHERE id=?")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database for changing cart to order: %v \n", err)
		panic(err.Error())
	}
	//update cart discount
	appDiscount, err := db.Prepare("UPDATE carts SET discount=? WHERE id=?")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database for discount: %v \n", err)
		panic(err.Error())
	}
	//update cart purchase date
	addCurDate, err := db.Prepare("UPDATE carts SET date_purchased=? WHERE id=?")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database for adding date for adding current date: %v \n", err)
		panic(err.Error())
	}
	//update customer balance
	upCustBalance, err := db.Prepare("UPDATE customers SET balance=balance-? WHERE id=?")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database for updating customer balance: %v \n", err)
		panic(err.Error())
	}
	//update customer discount date if amount is greater tha the given amount
	if totalPrice > config.ConfigInstance.GivenAmount {
		upCustDate, err := db.Prepare("UPDATE customers SET HAS_SUBSEQUENT_DISCOUNT_UNTIL=? WHERE id=?")
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error preparing database for updating customer balance: %v \n", err)
			panic(err.Error())
		}
		upCustDate.Exec(time.Now().Add(time.Hour*24*30), model.GetCartOwnerId(model.GetCart(cartId))) //update the customer discount date to 30 days from now
	}
	/* execute prepared statements */
	upCustBalance.Exec(totalPrice, model.GetCartOwnerId(model.GetCart(cartId)))
	addCurDate.Exec(time.Now(), cartId)
	appDiscount.Exec(discountAmount, cartId)
	completeForm.Exec(totalPrice, cartId)

	logger.AppLogger.Info().Printf("Cartid: %s, total price: %f, discount: %f succesfully completed.\n", cartId, totalPrice, discountAmount)

	//return response
	w.Write([]byte("Cart id: " + cartId + " completed."))
	w.Write([]byte("Total price: " + strconv.FormatFloat(totalPrice, 'f', 2, 64) + "Owner id: " + strconv.Itoa(model.GetCartOwnerId(model.GetCart(cartId)))))
	defer db.Close()

}

func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: DeleteCartItem")
	db := connector.DbConn()  //connect to the database
	if r.Method != "DELETE" { //check if the method is delete
		http.Error(w, http.StatusText(STATUS_NOT_ALLOWED), http.StatusMethodNotAllowed)
		w.Write([]byte("Method: DELETE is required."))
		logger.AppLogger.Error().Printf("Method Delete is required. \n")
		return
	}
	params := mux.Vars(r)             //get the params from the url
	cartId := params["cart_id"]       //get the cart id from the url
	productId := params["product_id"] //get the product id from the url

	if !model.IsCartExists(cartId) { //check if the cart exists
		http.Error(w, http.StatusText(STATUS_NOT_FOUND), http.StatusNotFound)
		w.Write([]byte("Cart: " + cartId + " not found"))
		logger.AppLogger.Error().Printf("Cart: %s not found", cartId)
		return
	}
	if !model.IsProductExists(productId) { //check if the product exists
		http.Error(w, http.StatusText(STATUS_NOT_FOUND), http.StatusNotFound)
		w.Write([]byte("Product: " + productId + " not found in the DB"))
		logger.AppLogger.Error().Printf("Product not found in DB: %s", productId)
		return
	}

	if !model.IsProductExistsInCart(cartId, productId) { //check if the product exists in the cart
		http.Error(w, http.StatusText(STATUS_NOT_FOUND), http.StatusNotFound)
		w.Write([]byte("Product: " + productId + " not found in cart: " + cartId))
		logger.AppLogger.Error().Printf("Product: %s not found in cart: %s", productId, cartId)
		return
	}

	logger.AppLogger.Info().Printf("Deleting product id: %s from cart id: %s \n", productId, cartId)
	//delete the product from the cart
	delForm, err := db.Prepare("DELETE FROM cart_items WHERE cart_id=? AND product_id=? LIMIT 1")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	delForm.Exec(cartId, productId)
	logger.AppLogger.Info().Println("Product deleted from cart.")
	model.UpdateCartPrice(cartId) //update the cart price

	w.Write([]byte("Product: " + productId + " deleted from cart: " + cartId))
	defer db.Close()
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: ShowCart")
	logger.AppLogger.Info().Println("Endpoint hit : ShowCart")
	db := connector.DbConn()         //connect to the database
	params := mux.Vars(r)            //get the params from the url
	cartId := params["cart_id"]      //get the cart id from the url
	if !model.IsCartExists(cartId) { //check if the cart exists
		logger.AppLogger.Error().Println("Cart not found")
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	//get the cart items from the cart
	selCart, err := db.Query("SELECT * FROM carts WHERE id = ?", cartId)
	if err != nil {
		logger.AppLogger.Error().Println(err.Error())
		panic(err.Error())
	}
	cart := model.Cart{}
	for selCart.Next() {
		/*we don't really need to use for since, we expect every cart to have a different id, but it'd be a good way to debug
		since we can check if the query above returns multiple carts
		*/
		var id int
		var customerId int
		var isPurchased bool
		var DatePurchased time.Time
		var TotalPrice float64
		var Discount float64

		err = selCart.Scan(&id, &customerId, &isPurchased, &DatePurchased, &TotalPrice, &Discount) //scan the cart data
		if err != nil {
			logger.AppLogger.Error().Println(err.Error())
			panic(err.Error())
		}
		//assign cart data to the cart struct
		cart.Id = id
		cart.CustomerId = customerId
		cart.IsPurchased = isPurchased
		cart.DatePurchased = DatePurchased
		cart.TotalPrice = TotalPrice
		cart.Discount = Discount
		logger.AppLogger.Info().Printf("Reading cart id: %d , customer id: %d , is purchased: %v , date purchased: %v , total price: %f \n", cart.Id, cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice)
	}
	json.NewEncoder(w).Encode(cart) //encode the cart data to json

	//get the cart items from the cart
	selDB, err := db.Query("SELECT * FROM cart_items WHERE cart_id = ?", cartId)
	if err != nil {
		logger.AppLogger.Error().Println(err.Error())
	}
	cartItem := model.CartItem{}    //create a cart item struct
	cartItems := []model.CartItem{} //create a cart item slice
	logger.AppLogger.Info().Println("Parsing cart items...")
	for selDB.Next() { //loop through the cart items
		var id int
		var cartId int
		var productId int
		var productName string

		err = selDB.Scan(&id, &cartId, &productId, &productName)
		if err != nil {
			logger.AppLogger.Error().Println(err.Error())
			panic(err.Error())
		}
		//assign cart item data to the cart item struct
		cartItem.Id = id
		cartItem.CartId = cartId
		cartItem.ProductId = productId
		cartItem.ProductName = productName

		cartItems = append(cartItems, cartItem) //append the cart item to the cart item slice
		logger.AppLogger.Info().Printf("Reading cart item id: %d , cart id: %d , product id: %d , product name: %s \n", cartItem.Id, cartItem.CartId, cartItem.ProductId, cartItem.ProductName)
	}
	json.NewEncoder(w).Encode(cartItems) //encode the cart item slice to json

	logger.AppLogger.Info().Printf("Total of %v cart items found and listed. \n", len(cartItems))
	defer db.Close()
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: AddToCart")
	db := connector.DbConn()
	if r.Method != "POST" { //check if the request is a POST request
		logger.AppLogger.Error().Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)                              //get the params from the url
	cartId := params["cart_id"]                        //get the cart id from the url
	productId := params["product_id"]                  //get the product id from the url
	productName := model.GetProductNameByID(productId) //get the product name from the product id

	if !model.IsCartExists(cartId) { //check if the cart exists
		logger.AppLogger.Error().Println("Cart not found")
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	if !model.IsProductExists(productId) { //check if the product exists
		logger.AppLogger.Error().Println("Product not found")
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	logger.AppLogger.Info().Printf("Adding product id: %s , product name: %s to cart id: %s \n", productId, productName, cartId)

	// insert the product to the cart
	insForm, err := db.Prepare("INSERT INTO cart_items (cart_id, product_id, product_name) VALUES(?,?,?)")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	insForm.Exec(cartId, productId, productName)
	w.Write([]byte("Product "+ productName +" added to cart "+ cartId+ " successfully")))) 
	logger.AppLogger.Info().Printf("Product id: %s , product name: %s added to cart id: %s \n", productId, productName, cartId)
	model.UpdateCartPrice(cartId) //update the cart price
	defer db.Close()
}

func ListCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: ListCustomers")
	logger.AppLogger.Info().Println("Endpoint hit : ListCustomers")
	db := connector.DbConn()
	//get all the customers from the database
	selDB, err := db.Query("SELECT * FROM customers")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	customer := model.Customer{}   //create a customer struct
	customers := []model.Customer{} //create a customer slice
	logger.AppLogger.Info().Println("Parsing customers...")
	for selDB.Next() { //loop through the customers
		var id int
		var name string
		var balance float64
		var consecutive_discount int
		var has_subsequent_discount_until time.Time
		err = selDB.Scan(&id, &name, &balance, &consecutive_discount, &has_subsequent_discount_until) //scan the customer data
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		//assign customer data to the customer struct
		customer.Id = id
		customer.Name = name
		customer.Balance = balance
		customer.Consecutive_discount = consecutive_discount
		customer.Has_subsequent_discount_until = has_subsequent_discount_until
		logger.AppLogger.Info().Printf("Reading customer id: %d , name: %s , balance: %f , consecutive discount: %d , has subsequent discount until: %v \n", customer.Id, customer.Name, customer.Balance, customer.Consecutive_discount, customer.Has_subsequent_discount_until)
		customers = append(customers, customer) //append the customer to the customer slice
	}
	json.NewEncoder(w).Encode(customers) //encode the customer slice to json
	logger.AppLogger.Info().Printf("Total of %v customers found and listed. \n", len(customers)) //print the total of customers found
	defer db.Close()
}

func ListCartItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: ListCartItems")
	logger.AppLogger.Info().Println("Endpoint hit : ListCartItems")
	db := connector.DbConn()
	//get all the cart items from the database
	selDB, err := db.Query("SELECT * FROM cart_items")
	if err != nil {
		logger.AppLogger.Error().Println(err.Error())
	} 
	cartItem := model.CartItem{}  //create a cart item struct
	cartItems := []model.CartItem{} //create a cart item slice
	logger.AppLogger.Info().Println("Parsing cart items...")
	for selDB.Next() {
		var id int
		var cartId int
		var productId int
		var productName string

		err = selDB.Scan(&id, &cartId, &productId, &productName) //scan the cart item data
		if err != nil {
			logger.AppLogger.Error().Println(err.Error())
		}
		//assign cart item data to the cart item struct
		cartItem.Id = id
		cartItem.CartId = cartId
		cartItem.ProductId = productId
		cartItem.ProductName = productName
		cartItems = append(cartItems, cartItem)
		logger.AppLogger.Info().Printf("Reading cart item id: %d , cart id: %d , product id: %d , product name: %s \n", cartItem.Id, cartItem.CartId, cartItem.ProductId, cartItem.ProductName)
	}
	json.NewEncoder(w).Encode(cartItems) //encode the cart item slice to json
	logger.AppLogger.Info().Printf("Total of %v cart items found and listed. \n", len(cartItems))
	defer db.Close()
}

func ListCarts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: ListCarts")
	logger.AppLogger.Info().Println("Endpoint hit : ListCarts")
	db := connector.DbConn()
	//get all the carts from the database
	selDB, err := db.Query("SELECT * FROM carts")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	cart := model.Cart{} //create a cart struct
	carts := []model.Cart{} //create a cart slice
	logger.AppLogger.Info().Println("Parsing carts...")
	for selDB.Next() {
		var id int
		var customer_id int
		var is_purchased bool
		var date_purchased time.Time
		var total_price float64
		var discount float64

		err = selDB.Scan(&id, &customer_id, &is_purchased, &date_purchased, &total_price, &discount) 	//scan the cart data
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		//assign cart data to the cart struct
		cart.Id = id
		cart.CustomerId = customer_id
		cart.IsPurchased = is_purchased
		cart.DatePurchased = date_purchased
		cart.TotalPrice = total_price
		carts = append(carts, cart) //append the cart to the cart slice
		logger.AppLogger.Info().Printf("Reading cart id: %d , customer id: %d , is purchased: %t , date purchased: %v , total price: %f , discount: %f \n", cart.Id, cart.CustomerId, cart.IsPurchased, cart.DatePurchased, cart.TotalPrice, cart.Discount)
	}
	json.NewEncoder(w).Encode(carts) //encode the cart slice to json
	logger.AppLogger.Info().Printf("Total of %v carts found and listed. \n", len(carts)) //print the total of carts found
	defer db.Close()

}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: ListProducts")
	logger.AppLogger.Info().Println("Endpoint hit : ListProducts")
	db := connector.DbConn()
	//get all the products from the database
	selDB, err := db.Query("SELECT * FROM products")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	product := model.Product{} //create a product struct
	products := []model.Product{} //create a product slice
	logger.AppLogger.Info().Println("Parsing products...")
	for selDB.Next() {
		var id int
		var name string
		var price float64
		var vat float64

		err = selDB.Scan(&id, &name, &price, &vat) //scan the product data
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
		logger.AppLogger.Info().Printf("Reading Product, ID: %v , name: %v , price: %v , vat: %v \n", id, name, price, vat)
		//assign product data to the product struct
		product.Id = id
		product.Name = name
		product.Price = price
		product.Vat = vat
		products = append(products, product) //append the product to the product slice
	}
	json.NewEncoder(w).Encode(products) //encode the product slice to json
	logger.AppLogger.Info().Printf("Total of %v products found and listed. \n", len(products))
	defer db.Close()
}
