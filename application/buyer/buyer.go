//Handles buyer functions - talks to templates and apiclient
//gets login information from server.go
package buyer

import (
	"fmt"
	"net/http"
	"projectGoLive/application/apiclient"
	"projectGoLive/application/config"
	"projectGoLive/application/email"
	"projectGoLive/application/server"
	"projectGoLive/application/user_db"
	"strings"
)

type buyerStruct struct {
	Buyername   string
	Operation   string
	Mainmessage []string
	Items       []apiclient.ItemsDetails
	CostPerItem []string
	Totalcost   string
}

var buyerCartll CartLinkedList

func init() {
	buyerCartll = CartLinkedList{Head: nil, Size: 0}
}

func BuyerHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := server.GetUser(w, req)

	buyerToTemplate := buyerStruct{}
	buyerToTemplate.Buyername = user.Username
	buyerToTemplate.Operation = "view"

	allSellerItems, ok := apiclient.GetItem("", "", user.IsBuyer)
	if !ok {
		config.Error.Println("Unable to connect to Database!")
		return
	}

	allSellerItems = removeCartItems(allSellerItems)

	buyerToTemplate.Items = allSellerItems
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "List of all items: ")

	if req.Method == http.MethodPost {
		addthisitem := req.FormValue("product_id")
		newquantity := req.FormValue("newquantity")

		if addthisitem != "" {
			// Add this item to linked list
			item := convStringtoSlice(addthisitem)
			item.Quantity = config.ConvertToInt(newquantity)

			// check if item exists in linked list
			iteminll, index, err := buyerCartll.SearchItemandSellerName(item.Item, item.Username)
			if err != nil {
				// item does not exist in linked list, can add as a new node
				err := buyerCartll.AddNode(item)
				if err != nil {
					buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Server Error. Try again!")
				} else {
					// redirect to main index
					http.Redirect(w, req, "/buyer/buyercart", http.StatusSeeOther)
					return
				}
			} else {
				// item exists, need to update the item
				newitem := item
				newitem.Quantity = item.Quantity + iteminll.Quantity
				err := buyerCartll.WriteAtIndex(index, newitem)
				if err != nil {
					buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Server Error. Try again!")
				} else {
					// redirect to main index
					http.Redirect(w, req, "/buyer/buyercart", http.StatusSeeOther)
					return
				}
			}
		}
	}
	//display available items on browser
	config.TPL.ExecuteTemplate(w, "buyertemplate.gohtml", buyerToTemplate)
}

func LookForItemHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := server.GetUser(w, req)
	var oneItemAllSellers []apiclient.ItemsDetails

	buyerToTemplate := buyerStruct{}
	buyerToTemplate.Buyername = user.Username
	buyerToTemplate.Operation = "finditem"

	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Please choose an item to search: ")

	// process form submission , when buyer clicks submit
	if req.Method == http.MethodPost {
		chosenitem := req.FormValue("fruit")
		addthisitem := req.FormValue("product_id")

		if chosenitem != "" {
			allSellerItems, ok := apiclient.GetItem("", "", user.IsBuyer)
			if !ok {
				buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Server Error. Try again!")
			}
			allSellerItems = removeCartItems(allSellerItems)

			for _, item := range allSellerItems {
				if item.Item == chosenitem {
					oneItemAllSellers = append(oneItemAllSellers, item)
				}
			}
			buyerToTemplate.Items = oneItemAllSellers
		} else if addthisitem != "" {
			// Add this item to linked list
			item := convStringtoSlice(addthisitem)
			newquantity := req.FormValue("newquantity")
			item.Quantity = config.ConvertToInt(newquantity)

			// check if item exists in linked list
			iteminll, index, err := buyerCartll.SearchItemandSellerName(item.Item, item.Username)
			if err != nil {
				// item does not exist in linked list, can add as a new node
				err := buyerCartll.AddNode(item)
				if err != nil {
					buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Server Error. Try again!")
				} else {
					// redirect to main index
					http.Redirect(w, req, "/buyer/buyercart", http.StatusSeeOther)
					return
				}
			} else {
				// item exists, need to update the item
				newitem := item
				newitem.Quantity = item.Quantity + iteminll.Quantity
				err := buyerCartll.WriteAtIndex(index, newitem)
				if err != nil {
					buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Server Error. Try again!")
				} else {
					// redirect to main index
					http.Redirect(w, req, "/buyer/buyercart", http.StatusSeeOther)
					return
				}
			}

		}
	}
	config.TPL.ExecuteTemplate(w, "buyertemplate.gohtml", buyerToTemplate)
}

func CartHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := server.GetUser(w, req)

	buyerToTemplate := buyerStruct{}
	buyerToTemplate.Buyername = user.Username
	buyerToTemplate.Operation = "viewcart"

	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Shopping Cart: ")

	_, allitems := buyerCartll.GetAllItems()
	buyerToTemplate.Items = allitems

	buyerToTemplate.Totalcost, buyerToTemplate.CostPerItem = computeTotalCost(allitems)

	// process form submission , when buyer clicks submit
	if req.Method == http.MethodPost {
		addmore := req.FormValue("add_more")
		reset := req.FormValue("reset")
		checkout := req.FormValue("checkout")

		if addmore != "" {
			http.Redirect(w, req, "/buyer", http.StatusSeeOther)
			return
		} else if reset != "" {
			// reset shopping cart linked list
			buyerCartll = CartLinkedList{Head: nil, Size: 0}
			http.Redirect(w, req, "/buyer", http.StatusSeeOther)
			return
		} else if checkout != "" {
			// perform checkout
			allSellerItems, ok := apiclient.GetItem("", "", user.IsBuyer)
			if !ok {
				config.Error.Println("Unable to connect to Database!")
				return
			}
			_, allitems := buyerCartll.GetAllItems()
			allok := true

			for _, item := range allitems {
				ok := updateDB(user.IsBuyer, allSellerItems, item)
				allok = allok && ok
				if !ok {
					buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Error while performing check out!")
					buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Try again")
				} else {
					// If ok, remove that item from the linked list
					_, index, err := buyerCartll.SearchItemandSellerName(item.Item, item.Username)
					if err != nil {
						config.Error.Println("Not able to find item in linked list")
						config.Error.Println(err)
					} else {
						_, err := buyerCartll.Remove(index)
						if err != nil {
							config.Error.Println("Not able to remove item from linked list")
							config.Error.Println(err)
						}
					}
				}
			}
			if allok {
				// Send invoice email to buyer and sellers
				email.Sendemail(user.Username, allitems)

				http.Redirect(w, req, "/buyer/checkoutsuccess", http.StatusSeeOther)
				return
			}
		}
	}

	config.TPL.ExecuteTemplate(w, "buyercart.gohtml", buyerToTemplate)
}

//---------------------------------------------------------------------------
// Functions to display profile of buyer
//---------------------------------------------------------------------------
func ShowProfile(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path)
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := server.GetUser(w, req)

	// Buyer details :
	buyerdetails, _ := user_db.GetARecord(config.DB, user.Username)
	buyerfullname := buyerdetails.Fullname
	buyeraddress := buyerdetails.Address
	buyeremail := buyerdetails.Email
	buyerphone := buyerdetails.Phone

	buyerToTemplate := buyerStruct{
		Buyername:   user.Username,
		Operation:   "profile",
		Mainmessage: nil,
		Items:       nil,
	}

	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Login name : "+user.Username)
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Full name : "+buyerfullname)
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Address : "+buyeraddress)
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Email : "+buyeremail)
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Phone : "+buyerphone)

	config.TPL.ExecuteTemplate(w, "buyertemplate.gohtml", buyerToTemplate)
}

func CheckoutSuccessHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := server.GetUser(w, req)
	buyerToTemplate := buyerStruct{}
	buyerToTemplate.Buyername = user.Username
	buyerToTemplate.Operation = "checkoutsuccess"

	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Checkout successful!")
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Invoice has been emailed to your registered email address.")
	buyerToTemplate.Mainmessage = append(buyerToTemplate.Mainmessage, "Please make payment during collection!")
	buyerCartll = CartLinkedList{Head: nil, Size: 0}

	config.TPL.ExecuteTemplate(w, "buyercart.gohtml", buyerToTemplate)
}

func computeTotalCost(allitems []apiclient.ItemsDetails) (totalcost string, costperitem []string) {
	total := 0.0
	temptotal := 0.0
	for _, item := range allitems {
		temptotal = item.Cost * float64(item.Quantity)
		costperitem = append(costperitem, fmt.Sprintf("%.2f", temptotal))
		total = total + temptotal
	}
	return fmt.Sprintf("%.2f", total), costperitem
}

func updateDB(isBuyer bool, allSellerItems []apiclient.ItemsDetails, oneCartItem apiclient.ItemsDetails) bool {
	tempItem := apiclient.ItemsDetails{}
	for _, item := range allSellerItems {
		if item.Item == oneCartItem.Item && item.Username == oneCartItem.Username {
			if item.Quantity > oneCartItem.Quantity {
				tempItem.Quantity = item.Quantity - oneCartItem.Quantity
				tempItem.Cost = item.Cost
				tempItem.Item = item.Item
				tempItem.Username = item.Username
				ok := apiclient.UpdateItem(item.Item, item.Username, isBuyer, tempItem)
				return ok
			} else {
				ok := apiclient.DeleteItem(item.Item, item.Username, isBuyer)
				return ok

			}
		}
	}
	return false
}

func removeCartItems(allSellerItems []apiclient.ItemsDetails) []apiclient.ItemsDetails {
	tempItem := apiclient.ItemsDetails{}
	_, allcartitems := buyerCartll.GetAllItems()

	for _, cartitem := range allcartitems {
		for index, item := range allSellerItems {
			if item.Item == cartitem.Item && item.Username == cartitem.Username {
				if item.Quantity > cartitem.Quantity {
					tempItem.Quantity = item.Quantity - cartitem.Quantity
					tempItem.Cost = item.Cost
					tempItem.Item = item.Item
					tempItem.Username = item.Username
					allSellerItems[index] = tempItem
				} else {
					allSellerItems = append(allSellerItems[0:index], allSellerItems[index+1:]...)
				}
				continue
			}
		}
	}
	return allSellerItems
}

func convStringtoSlice(iteminput string) apiclient.ItemsDetails {
	iteminput = strings.Replace(iteminput, "{", "", -1)
	iteminput = strings.Replace(iteminput, "}", "", -1)
	itemslice := strings.Split(iteminput, " ")
	item := apiclient.ItemsDetails{}
	item.Item = itemslice[0]
	item.Quantity = config.ConvertToInt(itemslice[1])
	item.Cost = config.ConvertToFloat(itemslice[2])
	item.Username = itemslice[3]
	return item
}
