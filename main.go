package main

import (
	logger "final_project/logger"
	sql "final_project/mysql"
	rest "final_project/rest_functions"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	logger.AppLogger.Info().Println("Main started")
	sql.CreateDb()
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	logger.AppLogger.Info().Println("Handling requests...")

	//Reading endpoints
	myRouter.HandleFunc("/listproducts", rest.ListProducts).Methods("GET")
	myRouter.HandleFunc("/listcarts", rest.ListCarts).Methods("GET")
	myRouter.HandleFunc("/listcartitems", rest.ListCartItems).Methods("GET")
	myRouter.HandleFunc("/listcustomers", rest.ListCustomers).Methods("GET")
	myRouter.HandleFunc("/showcart/{cart_id}", rest.ShowCart).Methods("GET")

	//Writing endpoints
	myRouter.HandleFunc("/addtocart/{cart_id}/{product_id}", rest.AddToCart).Methods("POST")
	myRouter.HandleFunc("/deletecartitem/{cart_id}/{product_id}", rest.DeleteCartItem).Methods("DELETE")
	myRouter.HandleFunc("/completecart/{cart_id}", rest.CompleteCart).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter)) //listen to 10000 port
}
