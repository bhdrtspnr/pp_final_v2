package config

import (
	"encoding/json"
	logger "final_project/logger"
	"io/ioutil"
)

type Config struct {
	GivenAmount                     float64 `json:"given_amount"`
	Point1VatDiscount               float64 `json:"point1_vat_discount"`
	Point8VatDiscount               float64 `json:"point8_vat_discount"`
	Point18VatDiscount              float64 `json:"point18_vat_discount"`
	SubsequentPurchaseDiscount      float64 `json:"subsequent_purchase_discount"`
	ThreeSubsequentPurchaseDiscount float64 `json:"three_subsequent_purchase_discount"`
	DbHost                          string  `json:"db_host"`
	DbPort                          string  `json:"db_port"`
	DbUser                          string  `json:"db_user"`
	DbPass                          string  `json:"db_password"`
	DbName                          string  `json:"db_name"`
	DbType                          string  `json:"db_type"`
	SqlLogs                         string  `json:"sql_logs"`
}

var ConfigInstance Config = GetConfig()

func GetConfig() Config {
	logger.AppLogger.Info().Println("Function hit : GetConfig")
	var config Config
	configFile, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error reading config file: %v \n", err)
		panic(err.Error())
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error unmarshalling config file: %v \n", err)
		panic(err.Error())
	}
	return config
}
