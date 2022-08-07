package model

import (
	logger "final_project/logger"
	connector "final_project/mysql"
)

type Product struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Vat   float64 `json:"vat"`
}

func GetProductName(id string) string {
	logger.AppLogger.Info().Println("Function hit : GetProductName")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT name FROM products WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var name string
	for selDB.Next() {
		err = selDB.Scan(&name)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	return name
}

func GetProductById(id string) Product {
	logger.AppLogger.Info().Println("Function hit : GetProduct")
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	var product Product
	for selDB.Next() {
		err = selDB.Scan(&product.Id, &product.Name, &product.Price, &product.Vat)
		if err != nil {
			logger.AppLogger.Fatal().Printf("Error scanning database: %v \n", err)
			panic(err.Error())
		}
	}
	return product
}

func IsProductExists(id string) bool {
	db := connector.DbConn()
	selDB, err := db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error querying database: %v \n", err)
		panic(err.Error())
	}
	if selDB.Next() {
		logger.AppLogger.Info().Printf("Product id: %s \n", id)
		return true
	}
	return false
}
