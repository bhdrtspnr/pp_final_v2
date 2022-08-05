
  # MySQL and DB
Mysql is required in local to run this project.
Please head to https://dev.mysql.com/downloads/installer/ download the installer and install MySQL server if you haven't already.

Program will come with 100 users and 100 products predefined. When you run the main function program will connect to the MySQL server and execute the pre defined SQL script located at /sql/create_db.sql file. You can modify the file if you'd like to add other products or users.

  # Testing
You can test certain functions with the test.http file, it has pre generated REST requests, I suggest you to install REST Client addon for VSCode, (https://marketplace.visualstudio.com/items?itemName=humao.rest-client), you can send requests and see their results pane by pane.

You can also use postman (https://www.postman.com/) to test the API.

  # Logging
Program comes in with a built in logger. Everytime the an action is executed program will generate the necessary logs at /logs/D-MMMM-YYYY.log file depending on the current date.
