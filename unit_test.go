package main

import (
	_ "final_project/config"
	discount "final_project/discount"
	logger "final_project/logger"
	"final_project/model"
	connector "final_project/mysql"
	_ "os"
	"testing"
)

/*
This is probably the worst way to do a unit test but it works.
Program literally wipes everything from the DB and inserts necessary values in each test.
The reason why I did this is I faced a couple of problems:
1) If program has not run before, the tables will be empty and the test will fail.
    Since program generates the necessary tables and inserts the data in the runtime by executing the script at sql/.
2) If a test fails, the program will not be able to run the next test.
3) Some tests are coupled with the others, which is creating dirty data and dependency on the previous test.


There are also some helper functions to create the necessary test data, I was not sure If I should
create another script in the sql/ folder to generate the data, it would probably be a better solution.

Go to the project directory and run the following to command to execute the tests:
go test -v

*/
func AddProducts() {

	db := connector.DbConn()
	delProd, err := db.Prepare("DELETE FROM products WHERE id in ('1001','1002','1003')")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delProd.Exec()
	logger.AppLogger.Info().Printf("Deleted products")
	defer db.Close()

	selDB, err := db.Prepare("INSERT INTO products(id, name, price, vat) VALUES (?,?,?,?)")

	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
		panic(err.Error())
	}
	selDB.Exec(1001, "wooden spoon", 100, 0.01)
	selDB.Exec(1002, "ice cream stick", 100, 0.08)
	selDB.Exec(1003, "zebra", 100, 0.18)

	logger.AppLogger.Info().Printf("Added products")
}

func AddCart() {
	db := connector.DbConn()
	defer db.Close()
	selDB, err := db.Prepare("INSERT INTO carts(id,customer_id, is_purchased, date_purchased, total_price, discount) VALUES (99999,99999, 0, '2019-01-01', 0, 0)")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}

	selDB.Exec()
	logger.AppLogger.Info().Printf("Added cart 1")

}

func AddCustomer() {
	db := connector.DbConn()
	delUser, err := db.Prepare("DELETE FROM customers WHERE id = '99999'")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delProd, err := db.Prepare("DELETE FROM products")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}

	delCartItems, err := db.Prepare("DELETE FROM cart_items")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delCart, err := db.Prepare("DELETE FROM carts")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delProd.Exec()
	delCartItems.Exec()
	delCart.Exec()
	delUser.Exec()
	logger.AppLogger.Info().Printf("Deleted customer")
	defer db.Close()
	selDB, err := db.Prepare("INSERT INTO customers(id, name, balance, consecutive_discount, has_subsequent_discount_until) VALUES ('99999', 'TEST_USER', 50000, 7, '9999-01-01')")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	selDB.Exec()
	logger.AppLogger.Info().Printf("Added customer 99999")
}

func AddTestWoodenSpoon() { //1001
	AddProducts()
	db := connector.DbConn()
	defer db.Close()
	delete, err := db.Prepare("DELETE FROM cart_items WHERE cart_id = 99999")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delete.Exec()
	logger.AppLogger.Info().Printf("Deleted cart items")

	selDB, err := db.Prepare("INSERT INTO cart_items(cart_id, product_id,product_name) VALUES (99999, 1001,'wooden spoon')")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	for i := 0; i < 10; i++ {
		selDB.Exec()
	}

}

func AddTestIceCreamStick() { //id = 1002
	AddProducts()
	db := connector.DbConn()
	defer db.Close()
	delete, err := db.Prepare("DELETE FROM cart_items WHERE cart_id = 99999")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delete.Exec()
	logger.AppLogger.Info().Printf("Deleted cart items")

	selDB, err := db.Prepare("INSERT INTO cart_items(cart_id, product_id,product_name) VALUES (99999, 1002,'ice cream stick')")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	for i := 0; i < 10; i++ {
		selDB.Exec()
	}

}

func AddTestZebra() { //id = 1003
	AddProducts()
	db := connector.DbConn()
	defer db.Close()
	delete, err := db.Prepare("DELETE FROM cart_items WHERE cart_id = 99999")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	delete.Exec()
	logger.AppLogger.Info().Printf("Deleted cart items")

	selDB, err := db.Prepare("INSERT INTO cart_items(cart_id, product_id,product_name) VALUES (99999, 1003,'zebra')")
	if err != nil {
		logger.AppLogger.Fatal().Printf("Error preparing database: %v \n", err)
	}
	for i := 0; i < 10; i++ {
		selDB.Exec()
	}

}

func TestCalculateConsecutivePurchaseDiscount(t *testing.T) {
	AddCustomer()
	AddCart()
	AddTestWoodenSpoon()

	want := 100.0
	model.UpdateCartPrice("99999")
	result := discount.CalculateConsecutivePurchaseDiscount("99999")
	if result != want {
		t.Errorf("Expected %f, got %f", want, result)
	}
	t.Logf("Expected %f, got %f", want, result)
}

func TestCalculateGivenAmountDiscountWith1Percent(t *testing.T) {
	AddCustomer()
	AddCart()
	AddTestWoodenSpoon()
	want := 0.0
	model.UpdateCartPrice("99999")
	result := discount.CalculateGivenAmountDiscount("99999")
	if result != want {
		t.Errorf("Expected %f, got %f", want, result)
	}
	t.Logf("Expected %f, got %f", want, result)

}

func TestCalculateGivenAmountDiscountWith8Percent(t *testing.T) {
	AddCustomer()
	AddCart()
	AddTestIceCreamStick()
	want := 100.0
	model.UpdateCartPrice("99999")
	result := discount.CalculateGivenAmountDiscount("99999")
	if result != want {
		t.Errorf("Expected %f, got %f", want, result)
	}
	t.Logf("Expected %f, got %f", want, result)

}

func TestCalculateGivenAmountDiscountWith18Percent(t *testing.T) {
	AddCustomer()
	AddCart()
	AddTestZebra()
	want := 150.0
	model.UpdateCartPrice("99999")
	result := discount.CalculateGivenAmountDiscount("99999")
	if result != want {
		t.Errorf("Expected %f, got %f", want, result)
	}
	t.Logf("Expected %f, got %f", want, result)

}

func TestCalculateThreeSubsequentPurchaseDiscount(t *testing.T) {
	AddCustomer()
	AddCart()
	AddTestWoodenSpoon()
	want := 24.0
	model.UpdateCartPrice("99999")
	result := discount.CalculateThreeSubsequentPurchaseDiscount("99999")
	if result != want {
		t.Errorf("Expected %f, got %f", want, result)
	}
	t.Logf("Expected %f, got %f", want, result)
}
