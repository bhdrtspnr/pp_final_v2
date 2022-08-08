package config

import (
	"encoding/json"
	logger "final_project/logger"
	"io/ioutil"
)

/*
As it was required in the assignment, given_amount should be modifyable by the user, so while doing that, I assumed
that other parameters like how much discount should be applied to consecutive purchases, and how much discount should be
so I made them all modifyable. Also I added DB connection parameters to the config struct, since it may change from user to user.
*/

//this section is complete copy pasta from: https://www.farsightsecurity.com/blog/txt-record/goconfig-20160523/

//create config struct for config.json
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

//create config struct for config.json
var ConfigInstance Config = GetConfig()

func GetConfig() Config {
	logger.AppLogger.Info().Println("Function hit : GetConfig")
	var config Config                                        //create config struct
	configFile, err := ioutil.ReadFile("config/config.json") //read config.json
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error reading config file: %v \n", err)
		panic(err.Error())
	}
	err = json.Unmarshal(configFile, &config) //unmarshal config.json to config struct
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error unmarshalling config file: %v \n", err)
		panic(err.Error())
	}
	return config //return config struct
}
