package mysql

import (
	"database/sql"
	logger "final_project/logger"
	"fmt"
	"io/ioutil"
	"strings"

	config "final_project/config"

	_ "github.com/go-sql-driver/mysql"
)

var dbDriver = config.ConfigInstance.DbType
var dbUser = config.ConfigInstance.DbUser
var dbPass = config.ConfigInstance.DbPass
var dbName = config.ConfigInstance.DbName

func createSchema() {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error creating database: %v \n", err)
		panic(err)
	}
	logger.AppLogger.Info().Printf("Database %v created \n", dbName)
}

func DbConn() (db *sql.DB) {

	logger.AppLogger.Info().Println("Connecting to database...")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName+"?parseTime=true")

	if err != nil {
		fmt.Println("error validating sql.Open arguments")
		logger.AppLogger.Fatal().Println(err)
		panic(err.Error())
	}
	return db
}

func CreateDb() {
	createSchema()
	file, err := ioutil.ReadFile("sql/create_db.sql")
	db := DbConn()
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error reading file: %v \n", err)
	}

	requests := strings.Split(string(file), ";")
	logger.AppLogger.Fatal().Println("Parsing requests...")

	for _, request := range requests {
		result, err := db.Exec(request)
		if err != nil {
			logger.AppLogger.Error().Printf("Error executing request: %v \n", err)
		} else {
			if config.ConfigInstance.SqlLogs == "true" {
				logger.AppLogger.Info().Printf("Request executed: %v \n", result)
			}
		}

	}
}
