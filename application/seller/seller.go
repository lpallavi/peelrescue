//Seller functions, talks to apiclient.go, server.go and templates
package seller

import (
	"net/http"
	"projectGoLive/application/apiclient"
	"projectGoLive/application/config"
	"projectGoLive/application/server"
	"projectGoLive/application/user_db"
)

type sellerStruct struct {
	Sellername  string
	Operation   string
	Mainmessage []string
	Selleritems []apiclient.ItemsDetails
}

//---------------------------------------------------------------------------
// Functions to display all items added by seller
//---------------------------------------------------------------------------
// This method is used to view the items added by seller
func SellerHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := server.GetUser(w, req)

	sellerMessage := sellerStruct{
		Sellername:  user.Username,
		Operation:   "view",
		Mainmessage: nil,
		Selleritems: nil,
	}
	sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "List of items added: ")

	if user.IsBuyer { // Not possible
		config.Trace.Printf("Incorrect login information! username: %s  is not a seller.", user.Username)
		config.Error.Printf("Incorrect login information! username: %s  is not a seller.", user.Username)
		return
	}

	si, ok := apiclient.GetItem("", user.Username, user.IsBuyer) // get all items for this seller only
	if !ok {
		config.Trace.Printf("Unable to get item data for seller %s \n", user.Username)
		config.Error.Printf("Unable to get item data for seller %s \n", user.Username)
		return
	}

	sellerMessage.Selleritems = si
	config.TPL.ExecuteTemplate(w, "sellertemplate.gohtml", sellerMessage)
}

//---------------------------------------------------------------------------
// Functions to add an item for seller
//---------------------------------------------------------------------------
func AddItemHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := server.GetUser(w, req)

	sellerMessage := sellerStruct{
		Sellername:  user.Username,
		Operation:   "add",
		Mainmessage: nil,
		Selleritems: nil,
	}

	// process form submission , when seller clicks submit
	if req.Method == http.MethodPost {

		var item apiclient.ItemsDetails
		fruitname := req.FormValue("fruit")
		quantity := req.FormValue("quantity")
		cost := req.FormValue("cost")

		if fruitname != "" {
			item.Item = fruitname
			item.Quantity = config.ConvertToInt(quantity)
			item.Cost = config.ConvertToFloat(cost)
			item.Username = user.Username

			ok := apiclient.AddItem(item.Item, item.Username, user.IsBuyer, item)
			if !ok {
				sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Unable to add item\nIf item already exists, please update item, else try again!")
			} else {
				http.Redirect(w, req, "/seller", http.StatusSeeOther)
				return
			}
		}
	}
	config.TPL.ExecuteTemplate(w, "sellertemplate.gohtml", sellerMessage)
}

//---------------------------------------------------------------------------
// Functions to update an item added by seller
//---------------------------------------------------------------------------
func UpdateItemHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := server.GetUser(w, req)

	sellerMessage := sellerStruct{
		Sellername:  user.Username,
		Operation:   "update",
		Mainmessage: nil,
		Selleritems: nil,
	}

	if req.Method == http.MethodPost {

		var item apiclient.ItemsDetails
		fruitname := req.FormValue("fruit")
		quantity := req.FormValue("quantity")
		cost := req.FormValue("cost")

		if fruitname != "" {
			item.Item = fruitname
			item.Quantity = config.ConvertToInt(quantity)
			item.Cost = config.ConvertToFloat(cost)
			item.Username = user.Username

			ok := apiclient.UpdateItem(item.Item, item.Username, user.IsBuyer, item)
			if !ok {
				sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Unable to update item, Item does not exist!")
			} else {
				http.Redirect(w, req, "/seller", http.StatusSeeOther)
				return
			}
		}
	}
	config.TPL.ExecuteTemplate(w, "sellertemplate.gohtml", sellerMessage)
}

//---------------------------------------------------------------------------
// Functions to delete an item added by seller
//---------------------------------------------------------------------------
func DeleteItemHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := server.GetUser(w, req)

	sellerMessage := sellerStruct{
		Sellername:  user.Username,
		Operation:   "delete",
		Mainmessage: nil,
		Selleritems: nil,
	}

	if req.Method == http.MethodPost {

		var item apiclient.ItemsDetails
		fruitname := req.FormValue("fruit")

		if fruitname != "" {
			item.Item = fruitname
			item.Username = user.Username

			ok := apiclient.DeleteItem(item.Item, item.Username, user.IsBuyer)
			if !ok {
				sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Unable to delete item, Item does not exist!")
			} else {
				http.Redirect(w, req, "/seller", http.StatusSeeOther)
				return
			}
		}
	}
	config.TPL.ExecuteTemplate(w, "sellertemplate.gohtml", sellerMessage)
}

//---------------------------------------------------------------------------
// Functions to display profile of seller
//---------------------------------------------------------------------------
func ShowProfile(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := server.GetUser(w, req)

	// Seller details :
	sellerdetails, _ := user_db.GetARecord(config.DB, user.Username)
	sellerfullname := sellerdetails.Fullname
	selleraddress := sellerdetails.Address
	selleremail := sellerdetails.Email
	sellerphone := sellerdetails.Phone

	sellerMessage := sellerStruct{
		Sellername:  user.Username,
		Operation:   "profile",
		Mainmessage: nil,
		Selleritems: nil,
	}

	sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Here are your profile details :")
	sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Full name :"+sellerfullname)
	sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Address :"+selleraddress)
	sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Email :"+selleremail)
	sellerMessage.Mainmessage = append(sellerMessage.Mainmessage, "Phone :"+sellerphone)

	config.TPL.ExecuteTemplate(w, "sellertemplate.gohtml", sellerMessage)
}
