
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
    "sql_logs": "false"
    }
  # MySQL and DB
Mysql is required in local to run this project.
Please head to https://dev.mysql.com/downloads/installer/ download the installer and install MySQL server if you haven't already.

Program will come with 100 users and 100 products predefined. When you run the main function program will connect to the MySQL server and execute the pre defined SQL script located at /sql/create_db.sql file. You can modify the file if you'd like to add other products or users.

Also while creating the initial mysql server please use the credidentials stated above at Config section or change the config.json in order to connect the mysql server with the appropriate creditentials.

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
Program comes in with a built in logger. Everytime the an action is executed program will generate the necessary logs at /logs/D-MMMM-YYYY.log file depending on the current date.
