package email

import (
	"fmt"
	"net/smtp"
	apiclient "projectGoLive/application/apiclient"
	config "projectGoLive/application/config"
	user_db "projectGoLive/application/user_db"
	"strconv"
	"time"
)

type SellerInfo struct {
	Fullname  string
	Address   string
	Phone     string
	Email     string
	CartItems []apiclient.ItemsDetails
}

func Sendemail(buyername string, cartItems []apiclient.ItemsDetails) bool {

	// Email configuration
	from := "peelrescue@gmail.com"
	password := "blueappleredorange"
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	timeNow := time.Now()
	tNow := timeNow.Format("2006-01-02 15:04:05")

	// Create a map to store cart items for each seller
	var sellerMap map[string]SellerInfo
	sellerMap = make(map[string]SellerInfo)

	// Buyer details :
	buyerdetails, _ := user_db.GetARecord(config.DB, buyername)
	buyerfullname := buyerdetails.Fullname
	buyeraddress := buyerdetails.Address
	buyeremail := buyerdetails.Email
	buyerphone := buyerdetails.Phone

	// Get all user records :
	userdetails, _ := user_db.GetRecords(config.DB)

	for _, cartitem := range cartItems {
		for _, user := range userdetails {

			if user.Username == cartitem.Username {
				tempCartSlice := sellerMap[cartitem.Username].CartItems
				tempCartSlice = append(tempCartSlice, cartitem)
				sellerfullname := user.Fullname
				selleraddress := user.Address
				sellerphone := user.Phone
				selleremail := user.Email

				sellerMap[cartitem.Username] = SellerInfo{sellerfullname, selleraddress, sellerphone, selleremail, tempCartSlice}
			}
		}
	}

	//------------------------
	// Sending email to buyer.

	to := []string{
		buyeremail,
	}
	subject := fmt.Sprintf("Subject: Invoice from Peel Rescue! Your order at %s\n", tNow)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><style>table {width:800px; border: 8px solid black;} th, td {border: 2px solid rgb(242, 248, 163);width:50px;height:40px;text-align: center;}</style><body>")
	body = body + fmt.Sprintf("<h1>Hello Buyer : %s </h1> ", buyerfullname)

	alltotalcost := 0.0
	// Send email invoice to buyer
	for _, value := range sellerMap {
		body = body + fmt.Sprintf("<br>")
		body = body + fmt.Sprintf("<h2>Items purchased from seller : %s </h2> ", value.Fullname)
		body = body + fmt.Sprintf("<p>Pick up address : %s </p> ", value.Address)
		body = body + fmt.Sprintf("<p>Phone : %s </p> ", value.Phone)
		body = body + fmt.Sprintf("<p>Email : %s </p> ", value.Email)

		body = body + fmt.Sprintf("<table>")
		body = body + fmt.Sprintf("<tr><th>Item name</th><th>Quantity (In kgs)</th><th>Cost (In SGD)</th></tr>")

		totalcost := 0.0
		for _, cartitem := range value.CartItems {
			body = body + fmt.Sprintf("<tr>")
			body = body + fmt.Sprintf("<td>%s</td>", cartitem.Item)
			body = body + fmt.Sprintf("<td>%s</td>", strconv.Itoa(cartitem.Quantity))
			body = body + fmt.Sprintf("<td>%s</td>", floattostr(cartitem.Cost))
			body = body + fmt.Sprintf("</tr>")
			totalcost = totalcost + (cartitem.Cost * float64(cartitem.Quantity))
		}
		alltotalcost = alltotalcost + totalcost
		body = body + fmt.Sprintf("</table>")
		body = body + fmt.Sprintf("<br>")
		body = body + fmt.Sprintf("<h2>Total Cost : $%.2f</h2>", totalcost)
		body = body + fmt.Sprintf("<br>")
	}
	body = body + fmt.Sprintf("<h2>Total Amount Spent for this order : $%.2f</h2>", alltotalcost)
	body = body + fmt.Sprintf("<br>")
	body = body + fmt.Sprintf("<h3>Thank you for using Peel Rescue. </h3>")
	body = body + fmt.Sprintf("<h4>Saving Earth one Peel at a time!</h4>")
	body = body + "</body></html>"

	// Sending email to buyer
	bytemessage := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, bytemessage)
	if err != nil {
		config.Error.Println(err)
		return false
	}
	config.Info.Println("Email Sent Successfully to !", buyeremail)

	//-----------------------------------
	// Send email invoice to each seller

	for _, value := range sellerMap {
		// Each sellers email address.
		to := []string{
			value.Email,
		}
		subject := fmt.Sprintf("Subject: Invoice from Peel Rescue! For Buyer: %s at %s\n", buyerfullname, tNow)
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

		body := fmt.Sprintf("<html><style>table {width:800px; border: 8px solid black;} th, td {border: 2px solid rgb(242, 248, 163);width:50px;height:40px;text-align: center;}</style><body>")
		body = body + fmt.Sprintf("<h1>Hello Seller : %s </h1> ", value.Fullname)
		body = body + fmt.Sprintf("<div>Items purchased by buyer : %s </div> ", buyerfullname)
		body = body + fmt.Sprintf("<p>Address for delivery : %s </p> ", buyeraddress)
		body = body + fmt.Sprintf("<p>Buyer Phone : %s </p> ", buyerphone)
		body = body + fmt.Sprintf("<p>Buyer Email : %s </p> ", buyeremail)

		body = body + fmt.Sprintf("<table>")
		body = body + fmt.Sprintf("<tr><th>Item name</th><th>Quantity (In kgs)</th><th>Cost (In SGD)</th></tr>")
		totalcost := 0.0
		for _, cartitem := range value.CartItems {
			body = body + fmt.Sprintf("<tr>")
			body = body + fmt.Sprintf("<td>%s</td>", cartitem.Item)
			body = body + fmt.Sprintf("<td>%s</td>", strconv.Itoa(cartitem.Quantity))
			body = body + fmt.Sprintf("<td>%s</td>", floattostr(cartitem.Cost))
			body = body + fmt.Sprintf("</tr>")
			totalcost = totalcost + (cartitem.Cost * float64(cartitem.Quantity))
		}
		body = body + fmt.Sprintf("</table>")
		body = body + fmt.Sprintf("<h2>Total Cost : $%.2f</h2>", totalcost)

		// Message.
		body = body + fmt.Sprintf("<h3>Thank you for using Peel Rescue. </h3>")
		body = body + fmt.Sprintf("<h4>Saving Earth one Peel at a time!</h4>")
		body = body + "</body></html>"

		// Sending email.
		bytemessage := []byte(subject + mime + body)

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, bytemessage)
		if err != nil {
			config.Error.Println(err)
			return false
		}
		config.Info.Println("Email Sent Successfully to !", value.Email)

	}
	return true
}

func floattostr(input float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input, 'f', -1, 64)
}
