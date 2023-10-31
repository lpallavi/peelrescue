package user_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type UserDetails struct {
	Username string
	Password string
	Fullname string
	Isbuyer  bool
	Phone    string
	Address  string
	Email    string
}

// Function to get all records from the MYSQL database.
// The function takes in the handle to the database.
// It returns all the info of all users as an array of type UserDetails.
// It returns true when retrieval of records from the database is successful.
// It returns false when there is any error encountered and retrieval of records is not successful.
func GetRecords(db *sql.DB) ([]UserDetails, bool) {
	var ud []UserDetails

	results, err := db.Query("SELECT * FROM `userdetails`")
	if err != nil {
		log.Println("Not able to get user details")
		log.Println(err)
	}
	defer results.Close()

	for results.Next() {
		// map this type to the record in the table
		var uinfo UserDetails
		err = results.Scan(&uinfo.Username, &uinfo.Password, &uinfo.Fullname, &uinfo.Isbuyer, &uinfo.Phone, &uinfo.Address, &uinfo.Email)
		if err != nil {
			log.Println("Unable to get records")
			log.Println(err)
			return ud, false
		}
		ud = append(ud, uinfo)
	}
	return ud, true
}

// Function to get one record from the MYSQL database.
// The function takes in the handle to the database.
// Its also takes in the name of the user to search for, of type string.
// It returns all the info of one user of type UserDetails.
// It returns true when retrieval of the record from the database is successful.
// It returns false when there is any error encountered, and retrieval of record is not successful.
func GetARecord(db *sql.DB, uname string) (UserDetails, bool) {
	var ud UserDetails
	query := fmt.Sprintf("SELECT * FROM `userdetails` WHERE Username='%s'", uname)
	results, err := db.Query(query)
	if err != nil {
		log.Println("Unable to find a record")
		log.Println(err)
		return ud, false
	}
	defer results.Close()

	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&ud.Username, &ud.Password, &ud.Fullname, &ud.Isbuyer, &ud.Phone, &ud.Address, &ud.Email)
		if err != nil {
			log.Println("Unable to get the record")
			log.Println(err)
			return ud, false
		}
	}
	if ud.Username == "" {
		return ud, false
	}
	return ud, true
}

// Function to insert one record into the MYSQL database.
// The function takes in the handle to the database.
// Its also takes in the user details to insert of type UserDetails.
// It returns true when the user is inserted into the database successfully.
// It returns false when there is any error encountered, and user is not inserted successfully.
func InsertRecord(db *sql.DB, ud UserDetails) bool {
	query := fmt.Sprintf("INSERT INTO `userdetails`(Username, Password, Fullname, Isbuyer, Phone, Address, Email) VALUES ('%s','%s','%s',%t,'%s','%s','%s')", ud.Username, ud.Password, ud.Fullname, ud.Isbuyer, ud.Phone, ud.Address, ud.Email)
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Unable to insert the record")
		log.Println(err)
		return false
	} else {
		return true
	}
}

// Function to update an existing record in the MYSQL database.
// The function takes in the handle to the database.
// It also takes in the name of the user to update of type string, and the new details of the user of type UserDetails.
// It returns true when the user details is updated in the database successfully.
// It returns false when there is any error encountered, and user details is not updated successfully.
func EditRecord(db *sql.DB, uname string, ud UserDetails) bool {
	query := fmt.Sprintf("UPDATE `userdetails` SET Username='%s', Password='%s', Firstname='%s', Isbuyer=%t, Phone='%s', Address='%s', Email='%s') WHERE Username='%s'", ud.Username, ud.Password, ud.Fullname, ud.Isbuyer, ud.Phone, ud.Address, ud.Email, uname)
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
// It takes in the name of the user to delete of type string.
// It returns true when the user is deleted from the database successfully.
// It returns false when there is any error encountered, and user is not deleted successfully.
func DeleteRecord(db *sql.DB, uname string) bool {
	query := fmt.Sprintf(
		"DELETE FROM `userdetails` WHERE Username='%s'", uname)
	_, err := db.Query(query)
	if err != nil {
		log.Println("Unable to delete the record")
		log.Println(err)
		return false
	}
	return true
}
