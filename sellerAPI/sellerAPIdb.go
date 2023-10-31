package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type ItemsDetails struct {
	Item     string  `json:"Item"`
	Quantity int     `json:"Quantity"`
	Cost     float64 `json:"Cost"`
	Username string  `json:"Username"`
}

//-----------------------------------------------------------------------
// Functions for seller
//-----------------------------------------------------------------------
// Function to get all records from the MYSQL database.
// The function takes in the handle to the database.
// It returns all the info of all courses as an array of type SellerDetails.
// It returns true when retrieval of records from the database is successful.
// It returns false when there is any error encountered and retrieval of records is not successful.
func GetRecordsSeller(db *sql.DB, SN string) ([]ItemsDetails, bool) {
	var sd []ItemsDetails
	query := fmt.Sprintf("SELECT * FROM sellerAPIdb.itemsdetails WHERE Username='%s';", SN)
	results, err := db.Query(query)
	if err != nil {
		log.Println("Not able to get seller details")
		log.Println(err)
	}

	for results.Next() {
		// map this type to the record in the table
		var si ItemsDetails
		err = results.Scan(&si.Item, &si.Quantity, &si.Cost, &si.Username)
		if err != nil {
			log.Println("Unable to get records")
			log.Println(err)
			return sd, false
		}
		sd = append(sd, si)
	}
	return sd, true
}

// Following functions need to be updated for sellerAPIdb

// Function to get one record from the MYSQL database.
// The function takes in the handle to the database.
// Its also takes in the name of the course to search for, of type string.
// It returns all the info of one course of type SellerDetails.
// It returns true when retrieval of the record from the database is successful.
// It returns false when there is any error encountered, and retrieval of record is not successful.
func GetARecordSeller(db *sql.DB, IN string, SN string) (ItemsDetails, bool) {
	var si ItemsDetails
	query := fmt.Sprintf("SELECT * FROM sellerAPIdb.itemsdetails WHERE Username='%s' AND Item='%s';", SN, IN)
	results, err := db.Query(query)
	if err != nil {
		log.Println("Unable to find a record")
		log.Println(err)
		return si, false
	}
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&si.Item, &si.Quantity, &si.Cost, &si.Username)
		if err != nil {
			log.Println("Unable to get the record")
			log.Println(err)
			return si, false
		}
	}
	if si.Item == "" || si.Username == "" {
		return si, false
	}
	return si, true
}

// Function to insert one record into the MYSQL database.
// The function takes in the handle to the database.
// Its also takes in the course to insert of type SellerDetails.
// It returns true when the course is inserted into the database successfully.
// It returns false when there is any error encountered, and course is not inserted successfully.
func InsertRecordSeller(db *sql.DB, sid ItemsDetails) bool {
	query := fmt.Sprintf("INSERT INTO `itemsdetails` (Item, Quantity, Cost, Username) VALUES ('%s',%d, %f,'%s');", sid.Item, sid.Quantity, sid.Cost, sid.Username)
	_, err := db.Query(query)

	if err != nil {
		log.Println("Unable to insert the record")
		log.Println(err)
		return false
	}
	return true
}

// Function to update an existing record in the MYSQL database.
// The function takes in the handle to the database.
// It also takes in the name of the course to update of type string, and the new details of the course of type SellerDetails.
// It returns true when the course is updated in the database successfully.
// It returns false when there is any error encountered, and course is not updated successfully.
func EditRecordSeller(db *sql.DB, IN string, SN string, sid ItemsDetails) bool {
	query := fmt.Sprintf("UPDATE `itemsdetails` SET Item='%s', Quantity= %d, Cost=%f, Username='%s' WHERE Item='%s' AND Username='%s';", sid.Item, sid.Quantity, sid.Cost, sid.Username, IN, SN)
	_, err := db.Query(query)
	if err != nil {
		log.Println("Unable to edit the record")
		log.Println(err)
		return false
	}
	return true
}

// Function to delete an existing record in the MYSQL database.
// The function takes in the handle to the database.
// It takes in the name of the course to delete of type string.
// It returns true when the course is deleted from the database successfully.
// It returns false when there is any error encountered, and course is not deleted successfully.
func DeleteRecordSeller(db *sql.DB, IN string, SN string) bool {
	query := fmt.Sprintf("DELETE FROM `itemsdetails` WHERE Item='%s' AND Username='%s'", IN, SN)
	_, err := db.Query(query)
	if err != nil {
		log.Println("Unable to delete the record")
		log.Println(err)
		return false
	}
	return true
}

//-----------------------------------------------------------------------
// Functions for buyer
//-----------------------------------------------------------------------
// Function to get all records from the MYSQL database.
// The function takes in the handle to the database.
// It returns all the info of all courses as an array of type SellerDetails.
// It returns true when retrieval of records from the database is successful.
// It returns false when there is any error encountered and retrieval of records is not successful.
func GetRecordsBuyer(db *sql.DB) ([]ItemsDetails, bool) {
	var sd []ItemsDetails
	query := fmt.Sprintf("SELECT * FROM sellerAPIdb.itemsdetails;")
	results, err := db.Query(query)
	if err != nil {
		log.Println("Not able to get seller details")
		log.Println(err)
	}

	for results.Next() {
		// map this type to the record in the table
		var si ItemsDetails
		err = results.Scan(&si.Item, &si.Quantity, &si.Cost, &si.Username)
		if err != nil {
			log.Println("Unable to get records")
			log.Println(err)
			return sd, false
		}
		sd = append(sd, si)
	}
	return sd, true
}
