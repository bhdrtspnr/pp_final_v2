  # What it does?
    Through this project, A basket service developed using REST API and GO.
  Customers are be able to purchase existing products. 
  
  The functions of this service are as follows;
  1. List Products
  - Users are be able to list all products.
  
  2. Add To Cart
  - Users can add their products to the basket and the total of the basket
  changes accordingly.

  3. Show Cart
  - Users can list the products they have added to their cart and total price and
  VAT of the cart.
  
  4. Delete Cart Item
  - Users can remove the products added from the cart. Notice removing an item
  may change discount.

  5. Complete Order
  - Users can create an order with the products they add to their cart. 

  Some business rules
  1. Products always have price and VAT (Value Added Tax, or KDV). VAT might be
  different for different products. Typical VAT percentage is %1, %8 and %18.
  
  2. There might be discount in following situations:
    a. Every fourth order whose total is more than given amount may have discount
    depending on products. Products whose VAT is %1 don’t have any discount
    but products whose VAT is %8 and %18 have discount of %10 and %15
    respectively.
    
    b. If there are more than 3 items of the same product, then fourth and
    subsequent ones would have %8 off.

    c. If the customer made purchase which is more than given amount in a month
    then all subsequent purchases should have %10 off.
    
    d. Only one discount can be applied at a time. Only the highest discount should
    be applied.


  # How to use?
 1) Mysql is required in local to run this project.
Please head to https://dev.mysql.com/downloads/installer/ download the installer and install MySQL server if you haven't already.
If you don't know how to install MySQL please watch this video: https://www.youtube.com/watch?v=GIRcpjg-3Eg (it took me a few hours to figure out I had to install MySQL server in order to use MySQL locally)

2)Go to the pp_final_v2/config/config.json, edit the file appropriately, you can find the explanations in the Configs sections. To run the program you need to set the db_user and db_password accordingly.

3)Either build an executable with go build or execute go run . via terminal. You can test the software with either given examples in example_requests.http, you can also craft your own requests with the following parameters:

  //Reading endpoints
	myRouter.HandleFunc("/showcart/{cart_id}", rest.ShowCart).Methods("GET")

	//Writing endpoints
	/addtocart/{cart_id}/{product_id} Methods("POST")
	/deletecartitem/{cart_id}/{product_id} Methods("DELETE")
	/completecart/{cart_id} Methods("POST")

 4)Program will come with 100 predefined users and products you can view them after running the program once or you can view them in the /sql/create_db.sql file.
 
 5)You can also test the business logic with the predefined unit tests in the unit_test.go file. 
 
 6)You can view logs in /logs directory with a log file named $current_day
 
 7)You can change config to change the business logic (discounts).

  # Discounts
  
  There are 3 types of discounts available:
1)  CalculateConsecutivePurchaseDiscount -> checks if the customer made a purchase more than the given amount in the config file in the last 30 days, if did, customer receives a discount amount of Config.SubsequentPurchaseDiscount (it is modifiable). (It can be done by setting the customer's HAS_SUBSEQUENT_DISCOUNT_UNTIL attribute to current date + 30 days, whenever he makes a purchase greater than the given amount, if current date is greater than the attribute, do not apply discount, if not apply discount.)

2) CalculateGivenAmountDiscount -> Checks if the customer's cart's total price is more than the config.GivenAmount (modifiable), if it is, checks if the customer had 3 other previous purchases with more than config.GivenAmount (we do this by incrementing the customer's CONSECUTIVE_DISCOUNT attribute by 1 each time he makes a purchase which has a greater cart value than given amount), if yes, it applies config.Point18VatDiscount to 18% VAT items, config.Point8VatDiscount to 8% VAT items and config.Point1VatDiscount to 1% VAT items. They're all modifiable since they may need to change in the future.
  
 3) CalculateThreeSubsequentPurchaseDiscount-> Checks if the cart has more than 3 of the same item. Like if customer purchases 4 apples and an apple costs 10$, customer pays 3*10$ for the first 3 apples and receives a discount on the 4th apple by %config.ThreeSubsequentPurchaseDiscount (modifiable). If customer purchases 8 apples, he/she receives the discount twice for the 4th and 8th item.
 
 4) Cart only utilizes the highest amount of the discounts avaiable.
  

  # Config
    {
    "given_amount": 100, //satisfies business logic 2.c and expectations 6
    "point1_vat_discount": 0, //satisfies business logic 2.a
    "point8_vat_discount": 0.1, //satisfies business logic 2.a
    "point18_vat_discount": 0.15, //satisfies business logic 2.a
    "subsequent_purchase_discount": 0.1, //satisfies business logic 2.c 
    "three_subsequent_purchase_discount": 0.08, //satisfies business logic 2.b
    "db_name": "app_db", 
    "db_user": "root",
    "db_password": "123456",
    "db_type": "mysql",
    "sql_logs": "false" //some sql statements create a lot of unnecessary logs, you can turn them back on by changing this to true
    }
  # MySQL and DB
Mysql is required in local to run this project.
Please head to https://dev.mysql.com/downloads/installer/ download the installer and install MySQL server if you haven't already.

Program will come with 100 users and 100 products predefined. When you run the main function program will connect to the MySQL server and execute the pre defined SQL script located at /sql/create_db.sql file. You can modify the file if you'd like to add other products or users.

Also while creating the initial mysql server please use the credidentials stated above at Config section or change the config.json in order to connect the mysql server with the appropriate creditentials.

Detailed view of the Database Model:
![image](https://user-images.githubusercontent.com/97244264/183481147-8026fa43-203b-417b-9b5b-1581a6ee92eb.png)


  # Testing
You can test certain functions with the test.http file, it has pre generated REST requests, I suggest you to install REST Client addon for VSCode, (https://marketplace.visualstudio.com/items?itemName=humao.rest-client), you can send requests and see their results pane by pane.

You can also use postman (https://www.postman.com/) to test the API.

Program comes with builtin unit tests as well, check unit_test.go file in the main section, you can either run all the tests at once and see the results with typing "go test -v" to terminal while in the main directory, or use VsCode IDE to run each test by clicking the button:

![image](https://user-images.githubusercontent.com/97244264/183280850-4524d711-2399-4c5c-b4ef-1659bc6e7221.png)

Notes about the unit tests (about how bad it is):
This is probably the worst way to do a unit test but it works.
Program literally wipes everything from the DB and inserts necessary values in each test.
The reason why I did this is I faced a couple of problems:
1) If program has not run before, the tables will be empty and the test will fail.
    Since program generates the necessary tables and inserts the data in the runtime by executing the script at sql/.
2) If a test fails, the program will not be able to run the next test.
3) Some tests are coupled with the others, which is creating dirty data and dependency on the previous test.


There are also some helper functions to create the necessary test data, I was not sure If I should
create another script in the sql/ folder to generate the data, it would probably be a better solution.


  # Logging
Program comes in with a built in logger. Everytime the an action is executed program will generate the necessary logs at /logs/D-MMMM-YYYY.log file depending on the current date. Logging was and is essential because without any notifications or information about what the program is doing, it's harder to debug or understand the logic working. You can pretty much understand what the program is doing when you execute a command.

As an example:

    POST http://localhost:10000/addtocart/1/10
     
    Generates the following output:
    INFO: 2022/08/07 10:52:15 config.go:28: Function hit : GetConfig
    INFO: 2022/08/07 10:52:15 main.go:14: Main started
    INFO: 2022/08/07 10:52:15 db_connector.go:32: Database app_db created 
    INFO: 2022/08/07 10:52:16 db_connector.go:37: Connecting to database...
    FATAL: 2022/08/07 10:52:16 db_connector.go:57: Parsing requests...
    ERROR: 2022/08/07 10:52:16 db_connector.go:62: Error executing request: Error 1065: Query was empty 
    INFO: 2022/08/07 10:52:16 main.go:21: Handling requests...
    INFO: 2022/08/07 10:52:21 db_connector.go:37: Connecting to database...
    INFO: 2022/08/07 10:52:21 product.go:16: Function hit : GetProductName
    INFO: 2022/08/07 10:52:21 db_connector.go:37: Connecting to database...
    INFO: 2022/08/07 10:52:21 cart.go:168: Function hit : IsCartExists
    INFO: 2022/08/07 10:52:21 db_connector.go:37: Connecting to database...
    INFO: 2022/08/07 10:52:21 db_connector.go:37: Connecting to database...
    INFO: 2022/08/07 10:52:21 product.go:61: Product id: 20 
    INFO: 2022/08/07 10:52:21 rest_functions.go:258: Adding product id: 20 , product name: basketball to cart id: 1 
    INFO: 2022/08/07 10:52:21 rest_functions.go:265: Product added to cart.
    INFO: 2022/08/07 10:52:21 cart.go:55: Function hit : UpdateCartPrice
    INFO: 2022/08/07 10:52:21 db_connector.go:37: Connecting to database...
    INFO: 2022/08/07 10:52:21 cart.go:110: Cart id: 1 price updated!, OLD VAL: 0.000000 , NEW_VAL: 21.000000 

