# peelrescue
Peel Rescue is an e-commerce shopping portal for buying/selling fruits peels to help reduce food wastage, and recover the abundance of nutrients present in fruit peel and put it to good use.
This project is built upon the concepts learnt at GoSchool during the period January - April 2021.
The theme of the project is to promote sustainability and reduce carbon footprint by choosing to recycle and reuse products such as a fruit peel that is normally just discarded without a second thought. This is one small step towards reducing the effects on our environment and climate change due to wastage of food. 
The primary users of this online website are fruit peel generators such as fruit juice shop owners who can sell leftover fruit peels to several fruit peel consumers which are industries that are looking for high value nutrients to be used in their products.

# Scope
Build an e-commerce portal to demonstrate the concepts learnt at GoSchool:
	Data structures:  Linked list for organizing cart items for buyer
	Microservice: Created a REST API with JSON for storing seller’s item information
	Security:  
o	SSL/TLS  Protocol for Application, and REST API
o	Login/Signup validation
o	Concurrent login detection
o	API keys
o	UUID for cookies
o	Hashed password, constant time verification of password
	Concurrency, mutex : Concurrency used for session management, and mutex for linked list
	Error handling: Reporting errors/warnings/info to log file
	Containerization and database : Docker  and mySQL
	Testing : Basic DB testing, and Client Console application for testing REST API
	Basic Knowledge of QA Test Cases: Third party package Data-Dog SQL mock has been used to mock up database functionalities

# Setup and info

1.	Copy the application folder “projectgolive” to a folder in PC’s Golang path e.g. D:/Golang/src/projectgolive

2.	To execute them main application : Execute the application using source code with the command 
1.	Change directory to application
2.	“go run main.go”
o	
3.	To execute them Seller REST API : Execute the API using source code with the command 
3.	Change directory to sellerAPI
4.	“go run .”

4.	Go source code files for main application: /application
5.	main.go – The main file for Peel Rescue to start the application
6.	config/* – This package declares all constants, connects to user DB, read admin credentials from env file, declares common functions, create template handle, and handle for error log file, and info/warning/trace logs
7.	start/* - This package declares all handler functions using gorilla mux router and Listen and Serve using TLS
8.	buyer/* – This package contains the linked list data structure for storing shopping cart and handlers for all buyer related functions.
9.	seller/* – This package contains all handlers for all seller related functions
10.	admin/* – This package contains all handlers for all admin related functions
11.	apiclient/* – This package contains all functions used to communicate with REST API using JSON
12.	server/* – This package contains all handler functions used for registration, user/admin login and logout. It also contains session management functions, cookie creation and deletion, input validation and sanitization. 
13.	user_db/* - This package contains all functions to interact with database storing user data
14.	templates/* - This folder contains all templates used in the main application
15.	assets/* - This folder contains all the assets used in the main application 	
16.	cert – This folder stores SSL certificate(cert.pem) and private key(key.pem)
17.	log – This folder stores log file for error log

5.	Go source code files for REST API: /sellerAPI
1.	sellerAPI.go – The file contains all functions to handler HTTP requests such as POST/GET/PUT AND DELETE
2.	sellerAPIdb.go – This file contains functions to interface with DB maintaining seller items


