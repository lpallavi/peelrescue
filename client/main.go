/*
Go file containing package main for running client console application.
The client application that allows the user to:
- Add Item
- Update Item
- Delete Item
- Retrieve Item
The package ignores TLS connection security as self-generated certificates are used for this project.
*/
package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

//Data structure for each item
type ItemsDetails struct {
	Item     string  `json:"Item"`
	Quantity int     `json:"Quantity"`
	Cost     float64 `json:"Cost"`
	Username string  `json:"Username"`
}

// Variable used only within this package
var sellerapikey string
var buyerapikey string
var apikey string

// Base URL used for testing server REST API
const baseURL = "https://localhost:5000/api/v1/"

// Initialization function
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// This function sends a request to the REST API to get one or all Items, and then displays the response.
// It ignores TLS security as REST API server uses self generated certicates
// It takes in the name of the Item to search, of type string
// If code is empty, it sends a request to search all Items
// Upon receiving the response from REST API, it displays the status of the request and the Item details.
func getItem(IN, SN string, isSeller bool) ([]ItemsDetails, bool) {
	var Items []ItemsDetails
	url := ""
	if isSeller {
		url = baseURL + "seller"
		apikey = sellerapikey
	} else {
		url = baseURL + "buyer"
		apikey = buyerapikey
	}

	if IN != "" {
		url = url + "/" + SN + "/" + IN
	} else if SN != "" {
		url = url + "/" + SN
	}

	// Skipping TLS verification as self generated certificate is used
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Get(url + "?key=" + apikey)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		if response.StatusCode == 200 {
			if IN != "" { // get one Item
				var oneItem ItemsDetails
				err := json.Unmarshal(data, &oneItem)
				if err != nil {
					log.Println(err)
				} else {
					Items = append(Items, oneItem)
					fmt.Println("Details of Item are : ")
					fmt.Printf("Item: \"%s\"\n", oneItem.Item)
					fmt.Printf("Quantity: %d\n", oneItem.Quantity)
					fmt.Printf("Cost: %f\n", oneItem.Cost)
					fmt.Printf("Username: \"%s\"\n", oneItem.Username)
					fmt.Println()
					// return one item in Items array, and true for successful get
					return Items, true
				}
			} else { // all Items
				err := json.Unmarshal(data, &Items)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println("List of all Items : ")
					for i, item := range Items {
						fmt.Printf("------- %d -------\n", i+1)
						fmt.Printf("Item: \"%s\"\n", item.Item)
						fmt.Printf("Quantity: %d\n", item.Quantity)
						fmt.Printf("Cost: %f\n", item.Cost)
						fmt.Printf("Username: \"%s\"\n", item.Username)
					}
					fmt.Println()
					// return all items in Items array, and true for successful get
					return Items, true
				}
			}
		} else if response.StatusCode == 404 {
			fmt.Println("Item not found. Try again")
			fmt.Println()
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(string(data))
			fmt.Println()
		}
	}
	// This return is for all errors, Items array will be empty, and false is not successful
	return Items, false
}

// This function sends a request to the REST API to add one Item, and then displays the response.
// It ignores TLS security as REST API server uses self generated certicates
// It takes in the name of the Item to add of type string.
// It also takes in the json data to be sent containing details of the Item to add.
// Upon receiving the response from REST API, it displays the status of the request and if Item has been added successfully.
//func addItem(code string, jsonData map[string]string) {
func addItem(IN, SN string, isSeller bool, si ItemsDetails) bool {
	url := ""
	if isSeller {
		apikey = sellerapikey
		url = baseURL + "seller"
	} else {
		return false
	}

	jsonValue, _ := json.Marshal(si)

	// Skipping TLS verification as self generated certificate is used
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Post(url+"/"+SN+"/"+IN+"?key="+apikey,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		if response.StatusCode == 201 {
			fmt.Println("Item added successfully.")
			fmt.Println()
			return true
		} else if response.StatusCode == 409 {
			fmt.Println("Item already exists! Try again.")
			fmt.Println()
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(string(data))
			fmt.Println()
		}
	}

	return false
}

// This function sends a request to the REST API to update one Item, and then displays the response.
// It ignores TLS security as REST API server uses self generated certicates
// It takes in the name of the Item to update of type string.
// It also takes in the json data to be sent containing details of the Item to update.
// Upon receiving the response from REST API, it displays the status of the request and if Item has been updated successfully.
//func updateItem(code string, jsonData map[string]string) {
func updateItem(IN, SN string, isSeller bool, si ItemsDetails) bool {
	url := ""
	if isSeller {
		apikey = sellerapikey
		url = baseURL + "seller"
	} else {
		apikey = buyerapikey
		url = baseURL + "buyer"
	}

	jsonValue, _ := json.Marshal(si)

	request, err := http.NewRequest(http.MethodPut,
		url+"/"+SN+"/"+IN+"?key="+apikey,
		bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	// Skipping TLS verification as self generated certificate is used
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		if response.StatusCode == 201 {
			fmt.Println("Item not in database. Added as a new Item.")
			fmt.Println()
			return true
		} else if response.StatusCode == 202 {
			fmt.Println("Item updated successfully.")
			fmt.Println()
			return true
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(string(data))
			fmt.Println()
		}
	}
	return false
}

// This function sends a request to the REST API to delete one Item, and then displays the response.
// It ignores TLS security as REST API server uses self generated certicates
// It takes in the name of the Item to be deleted of type string.
// Upon receiving the response from REST API, it displays the status of the request and if Item has been deleted successfully.
func deleteItem(IN, SN string, isSeller bool) bool {
	url := ""
	if isSeller {
		apikey = sellerapikey
		url = baseURL + "seller"
	} else {
		apikey = buyerapikey
		url = baseURL + "buyer"
	}
	fmt.Println("URL is ", url)

	request, err := http.NewRequest(http.MethodDelete, url+"/"+SN+"/"+IN+"?key="+apikey, nil)

	// Skipping TLS verification as self generated certificate is used
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		if response.StatusCode == 202 {
			fmt.Println("Item deleted successfully.")
			fmt.Println()
			return true
		} else if response.StatusCode == 404 {
			fmt.Println("Item not found. Try again")
			fmt.Println()
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(string(data))
			fmt.Println()
		}
	}
	return false
}

// Function to validate title/instructor name of Item
// If returns true if title/instructor name is between 2 and 32 characters long and contains only alphabets, numbers, spaces or _
// First character must be an alphabet, remaining characters must be alphabet,digits or underscore
func isValidString(t string) bool {
	unMatch := regexp.MustCompile(`^[A-Za-z]+[\w ]+$`)
	if unMatch.MatchString(t) {
		if len(t) < 1 || len(t) > 32 {
			fmt.Println("Length of Item must be more than 2 letters, and less than 32 letters")
			return false
		} else {
			return true
		}
	} else {
		fmt.Println("First character must be an alphabet, remaining characters must be alphabets,digits,spaces or underscore")
		return false
	}
}

// Function to validate quantity of Item
// If returns true if quantity is above 0
func isValidInt(d int) bool {
	if d < 0 {
		return false
	} else {
		return true
	}
}

// Function to scan a string using bufio, that can contain spaces
// If returns a string that has been scanned, and error if any
func scanWithSpaces() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); !ok {
		return "", scanner.Err()
	} else {
		line := scanner.Text()
		rl := fmt.Sprintf("%s", line)
		return rl, nil
	}
}

// Function to get an input from user that is a valid string, and 2 to 32 characters long
// If returns the string input
func getStringInput() string {
	for {
		name, err := scanWithSpaces()
		if err != nil {
			fmt.Println("Invalid input..Try again")
		} else {
			return name
		}
	}
}

// Function to get an input from user that is a valid int
// If returns the integer input
func getIntInput() int {
	var ii int
	for {
		_, err := fmt.Scanln(&ii)
		if err != nil || !isValidInt(ii) {
			fmt.Println("Invalid input. Quantity must be non-negative. Try again")
		} else {
			return ii
		}
	}
}

// Function to get an input from user that is a valid float
// If returns the float input
func getFloatInput() float64 {
	var fi float64
	for {
		_, err := fmt.Scanln(&fi)
		if err != nil || fi < 0.0 {
			fmt.Println("Invalid input. Cost must not be negative. Try again")
		} else {
			return fi
		}
	}
}

// Function containing the console interface for the client application
// All user options are displayed using fmt package
// Case of title of Item is ignored when comparing Item name with mySQL database
// Spaces in title are not ignored
func client_console() {
	var choice int
	var SN, IN string
	var CT float64
	var QT int
	fmt.Println("Welcome to Console Application for performing CRUD operations.")
	for {
		fmt.Println("Please choose an option: ")
		fmt.Println("----SELLER-----")
		fmt.Println("1. View All Items for a seller")
		fmt.Println("2. View A Particular Item for a seller")
		fmt.Println("3. Add An item for a seller")
		fmt.Println("4. Update An item for a seller")
		fmt.Println("5. Delete An item for a seller")
		fmt.Println("----BUYER-----")
		fmt.Println("6. View All Items for a buyer")
		fmt.Println("7. View A Particular Item for a buyer")
		fmt.Println("8. Add An item for a buyer")
		fmt.Println("9. Update An item for a buyer")
		fmt.Println("10. Delete An item for a buyer")
		fmt.Println("11. Exit")
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			fmt.Println("FOR SELLER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			getItem("", SN, true) // get all Items
		case 2:
			fmt.Println("FOR SELLER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name:")
			IN = getStringInput()
			getItem(IN, SN, true) // get all Items
		case 3:
			fmt.Println("FOR SELLER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name:")
			IN = getStringInput()
			fmt.Println("Enter the Quantity:")
			QT = getIntInput()
			fmt.Println("Enter the Cost:")
			CT = getFloatInput()

			var si ItemsDetails
			si.Item = IN
			si.Quantity = QT
			si.Cost = CT
			si.Username = SN

			addItem(IN, SN, true, si)
		case 4:
			fmt.Println("FOR SELLER:")
			var si ItemsDetails
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name to update:")
			IN = getStringInput()

			fmt.Println("----DETAILS OF UPDATED Item----")
			fmt.Println("Enter the seller name:")
			si.Username = getStringInput()
			fmt.Println("Enter the Item name:")
			si.Item = getStringInput()
			fmt.Println("Enter the Quantity:")
			si.Quantity = getIntInput()
			fmt.Println("Enter the Cost:")
			si.Cost = getFloatInput()
			updateItem(IN, SN, true, si)

		case 5:
			fmt.Println("FOR SELLER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name to delete:")
			IN = getStringInput()
			deleteItem(IN, SN, true)
		case 6:
			fmt.Println("FOR BUYER:")
			getItem("", "", false) // get all Items for buyer
		case 7:
			fmt.Println("FOR BUYER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name:")
			IN = getStringInput()
			getItem(IN, SN, false) // get one item of seller for buyer
		case 8:
			fmt.Println("FOR BUYER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name:")
			IN = getStringInput()
			fmt.Println("Enter the Quantity:")
			QT = getIntInput()
			fmt.Println("Enter the Cost:")
			CT = getFloatInput()

			var si ItemsDetails
			si.Item = IN
			si.Quantity = QT
			si.Cost = CT
			si.Username = SN

			addItem(IN, SN, false, si)
		case 9:
			fmt.Println("FOR BUYER:")
			var si ItemsDetails
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name to update:")
			IN = getStringInput()

			fmt.Println("----DETAILS OF UPDATED Item----")
			fmt.Println("Enter the seller name:")
			si.Username = getStringInput()
			fmt.Println("Enter the Item name:")
			si.Item = getStringInput()
			fmt.Println("Enter the Quantity:")
			si.Quantity = getIntInput()
			fmt.Println("Enter the Cost:")
			si.Cost = getFloatInput()
			updateItem(IN, SN, false, si)

		case 10:
			fmt.Println("FOR BUYER:")
			fmt.Println("Enter the seller name:")
			SN = getStringInput()
			fmt.Println("Enter the Item name to delete:")
			IN = getStringInput()
			deleteItem(IN, SN, false)

		case 11:
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid input..Try again")
		}
	}
}

// Function main for client application
func main() {
	/// Get the SELLER_API_KEY environment variable
	sellerapikey, _ = os.LookupEnv("SELLER_API_KEY")

	// Get the SELLER_API_KEY environment variable
	buyerapikey, _ = os.LookupEnv("BUYER_API_KEY")

	// Launch the client console application
	client_console()
}
