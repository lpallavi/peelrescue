/*
Go file containing package main for running REST API and functions for CRUD operations and database operations
The REST API connects to a MySQL database that stores information about different items.
The MYSQL database is deployed on a docker container using a MySQL database image.
The REST API allows items to be:
- Created
- Updated
- Deleted
- Retrieved
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Database handle for connecting to a MYSQL database
var sdb *sql.DB

const (
	// Directory that stores self generated cerificate.
	certPath = "./cert/"
)

// Variable used only within this package
var sellerapikey string
var buyerapikey string

//Initialization function
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Get the SELLER_API_KEY environment variable
	sellerapikey, _ = os.LookupEnv("SELLER_API_KEY")

	// Get the SELLER_API_KEY environment variable
	buyerapikey, _ = os.LookupEnv("BUYER_API_KEY")
}

// Function to validate the API KEY.
// It takes in the http response writer and http request as input and also the api key for either buyer or seller.
// It extracts the query in the request and checks if the api key in the url matches with the expected api key.
// It returns true is api keys match, and false if they dont.
// TO be updated for api
func validKey(w http.ResponseWriter, r *http.Request, apikey string) bool {
	v := r.URL.Query()

	if key, ok := v["key"]; ok {
		if key[0] == apikey {
			return true
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("401 - Invalid key"))
			return false
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Please supply access key"))
		return false
	}
}

// Function to display home page of REST API.
func sellerapihome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Seller API")
}

// Function to generate items page /api/v1/seller/{sellername}, containing all items information for one seller embedded in JSON format.
func seller_allitems(w http.ResponseWriter, r *http.Request) {

	if !validKey(w, r, sellerapikey) {
		log.Println("Seller API key not valid")
		return
	}
	params := mux.Vars(r)
	SN := params["sellername"]
	sid, ok := GetRecordsSeller(sdb, SN)
	if ok {
		// returns all the items in JSON
		json.NewEncoder(w).Encode(sid)
	}
}

// Function to generate items page for /api/v1/seller/{sellername}/{itemname}
// It handles GET/POST/PUT/DELETE methods sent from main application
// It generates all item information in response to request made, embedded in JSON format.
// It also generates headers for each request, depending on the status of each operation.
func seller_edititems(w http.ResponseWriter, r *http.Request) {
	if !validKey(w, r, sellerapikey) {
		log.Println("Seller API key not valid")
		return
	}
	params := mux.Vars(r)

	SN := params["sellername"]
	IN := params["itemname"]

	if r.Method == "GET" { // need not been checked if JSON or not
		item, ok := GetARecordSeller(sdb, IN, SN)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No items found"))
		} else {
			json.NewEncoder(w).Encode(item)
		}
	}

	if r.Method == "DELETE" {
		_, ok := GetARecordSeller(sdb, IN, SN)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No item found"))
		} else {
			DeleteRecordSeller(sdb, IN, SN)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Item deleted: " + IN + " For seller: " + SN))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating a new item
		if r.Method == "POST" {
			// read the string sent to the service
			var sid ItemsDetails
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply item information in JSON format"))
			} else {
				// convert JSON to object
				json.Unmarshal(reqBody, &sid)
				if sid.Username == "" || sid.Username != SN || sid.Item == "" || sid.Cost < 0.0 || sid.Quantity < 0 || sid.Item != IN {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply correct item information in JSON format"))
					return
				}
				// check if item exists; add only if the item does not exist
				_, ok := GetARecordSeller(sdb, IN, SN)
				if ok { // ok means its exists
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate item ID"))
				} else {
					InsertRecordSeller(sdb, sid)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Item added: " + IN + " For seller: " + SN))
				}
			}
		}

		//PUT is for creating or updating an existing item
		if r.Method == "PUT" {
			var sid ItemsDetails
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &sid)
				if sid.Username == "" || sid.Username != SN || sid.Item == "" || sid.Cost < 0.0 || sid.Quantity < 0 {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply item information in JSON format"))
					return
				}
				// check if item exists; add only if item does not exist
				_, ok := GetARecordSeller(sdb, IN, SN)
				if !ok { // item does not exist in db
					// Now check if new item does not exist either
					_, ok := GetARecordSeller(sdb, sid.Item, sid.Username)
					if ok { // exists, that means item for seller name does exist, but new item exists in db
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply correct item information in JSON format"))
						return
					} else { // item does not exist at all, need to add as new item
						InsertRecordSeller(sdb, sid)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Item added: " + IN + " For seller: " + SN))
					}
				} else {
					// update the item if item exists
					EditRecordSeller(sdb, IN, SN, sid)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Item updated: From " + IN + " To " + sid.Item + " For seller: " + SN))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply item information in JSON format"))
			}
		}
	}
}

// Function to generate items page /api/v1/buyer, containing all items information from all sellers embedded in JSON format.
func buyer_allitems(w http.ResponseWriter, r *http.Request) {
	if !validKey(w, r, buyerapikey) {
		log.Println("Buyer API key not valid")
		return
	}

	bid, ok := GetRecordsBuyer(sdb)
	if ok {
		// returns all the items in JSON
		json.NewEncoder(w).Encode(bid)
	}
}

// Function to generate items page for /api/v1/buyer/{sellername}/{itemname}
// It handles GET/PUT/DELETE methods sent from main application, POST cannot be done by buyer
// It generates all item information in response to request made, embedded in JSON format.
// It also generates headers for each request, depending on the status of each operation.
func buyer_edititems(w http.ResponseWriter, r *http.Request) {
	if !validKey(w, r, buyerapikey) {
		log.Println("Buyer API key not valid")
		return
	}
	params := mux.Vars(r)

	SN := params["sellername"]
	IN := params["itemname"]

	if r.Method == "GET" { // need not been checked if JSON or not
		item, ok := GetARecordSeller(sdb, IN, SN)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No items found"))
		} else {
			json.NewEncoder(w).Encode(item)
		}
	}

	if r.Method == "DELETE" {
		_, ok := GetARecordSeller(sdb, IN, SN)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No item found"))
		} else {
			DeleteRecordSeller(sdb, IN, SN)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Item deleted: " + IN + " For seller: " + SN))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating a new item
		if r.Method == "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 - Method Not Allowed"))
		}

		//PUT is for creating or updating an existing item
		if r.Method == "PUT" {
			var sid ItemsDetails
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &sid)
				if sid.Username == "" || sid.Username != SN || sid.Item == "" || sid.Cost < 0.0 || sid.Quantity < 0 {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply item information in JSON format"))
					return
				}
				// check if item exists; add only if item does not exist
				_, ok := GetARecordSeller(sdb, IN, SN)
				if !ok { // item does not exist in db
					// Now check if new item does not exist either
					_, ok := GetARecordSeller(sdb, sid.Item, sid.Username)
					if ok { // exists, that means item for seller name does exist, but new item exists in db
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply correct item information in JSON format"))
						return
					} else { // item does not exist at all, it cannot be added by buyer
						w.WriteHeader(http.StatusMethodNotAllowed)
						w.Write([]byte("405 - Method Not Allowed"))
					}
				} else {
					// update the item if item exists
					EditRecordSeller(sdb, IN, SN, sid)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Item updated: From " + IN + " To " + sid.Item + " For seller: " + SN))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply item information in JSON format"))
			}
		}
	}
}

// Function Main to open database, instantiate mux router, handle router functions, and listen and serve
func main() {

	// Open the mysql database seller_db created as a container in docker
	var err error
	sdb, err = sql.Open("mysql", "root:password@tcp(localhost:33061)/sellerAPIdb")
	// panic if unable to open database
	if err != nil {
		log.Println("Unable to open database")
		panic(err.Error())
	} else {
		log.Println("Database opened")
	}

	if err = sdb.Ping(); err != nil {
		panic(err)
	}

	// defer the database closing till after the main function has finished executing
	defer sdb.Close()

	// Instantiate mux router for handling urls
	router := mux.NewRouter()

	// API Home
	router.HandleFunc("/api/v1/", sellerapihome)

	// Handle function for all router functions for seller
	router.HandleFunc("/api/v1/seller/{sellername}", seller_allitems)                                                     // GET all items for a particular seller {sellername}
	router.HandleFunc("/api/v1/seller/{sellername}/{itemname}", seller_edititems).Methods("GET", "PUT", "POST", "DELETE") // one specific item for a particular seller {sellername}

	// Handle function for all router functions for buyer
	router.HandleFunc("/api/v1/buyer", buyer_allitems)                                                          // GET all items from all sellers
	router.HandleFunc("/api/v1/buyer/{sellername}/{itemname}", buyer_edititems).Methods("GET", "PUT", "DELETE") // one specific item from {sellername}, no POST

	// Listen and Serve TLS, using self generated cert.pem and key.pem
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServeTLS(":5000", certPath+"cert.pem", certPath+"key.pem", router))
}
