package model

import (
	logger "final_project/logger"
	connector "final_project/mysql"
	"time"
)

//create customer structure for customers table
type Customer struct {
	Id                            int       `json:"id"`
	Name                          string    `json:"name"`
	Balance                       float64   `json:"balance"`
	Consecutive_discount          int       `json:"consecutive_discount"`
	Has_subsequent_discount_until time.Time `json:"has_subsequent_discount_until"`
}

func GetCustomer(id int) Customer {
	//get customer by id
	logger.AppLogger.Info().Println("Function hit : GetCustomer")
	db := connector.DbConn()
	defer db.Close()
	//query the database for the customer with the given id
	selDB, err := db.Query("SELECT * FROM customers WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var customer Customer //create a customer object
	for selDB.Next() {
		err = selDB.Scan(&customer.Id, &customer.Name, &customer.Balance, &customer.Consecutive_discount, &customer.Has_subsequent_discount_until)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	return customer //return the customer object
}

func IsCustomerExist(customer Customer) bool {
	//check if customer exists in the database
	logger.AppLogger.Info().Println("Function hit : IsCustomerExist")
	db := connector.DbConn()
	//query the database for the customer with the given id
	selDB, err := db.Query("SELECT * FROM customers WHERE id = ?", customer.Id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() {
		logger.AppLogger.Info().Printf("Customer id: %d , name: %s , balance: %f , consecutive discount: %d , has subsequent discount until: %v \n", customer.Id, customer.Name, customer.Balance, customer.Consecutive_discount, customer.Has_subsequent_discount_until)
		return true //return true if customer exists
	}
	return false //return false if customer does not exist
}
