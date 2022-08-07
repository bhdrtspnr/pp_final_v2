package model

import (
	logger "final_project/logger"
	connector "final_project/mysql"
	"time"
)

type Customer struct {
	Id                            int       `json:"id"`
	Name                          string    `json:"name"`
	Balance                       float64   `json:"balance"`
	Consecutive_discount          int       `json:"consecutive_discount"`
	Has_subsequent_discount_until time.Time `json:"has_subsequent_discount_until"`
}

func GetCustomer(id int) Customer {
	logger.AppLogger.Info().Println("Function hit : GetCustomer")
	db := connector.DbConn()
	defer db.Close()
	selDB, err := db.Query("SELECT * FROM customers WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var customer Customer
	for selDB.Next() {
		err = selDB.Scan(&customer.Id, &customer.Name, &customer.Balance, &customer.Consecutive_discount, &customer.Has_subsequent_discount_until)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	return customer
}

func IsCustomerExist(customer Customer) bool {
	logger.AppLogger.Info().Println("Function hit : IsCustomerExist")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM customers WHERE id = ?", customer.Id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() {
		logger.AppLogger.Info().Printf("Customer id: %d , name: %s , balance: %f , consecutive discount: %d , has subsequent discount until: %v \n", customer.Id, customer.Name, customer.Balance, customer.Consecutive_discount, customer.Has_subsequent_discount_until)
		return true
	}
	return false
}
