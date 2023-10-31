package start

import (
	"net/http"
	"projectGoLive/application/admin"
	"projectGoLive/application/buyer"
	"projectGoLive/application/seller"
	"projectGoLive/application/server"
)

func mapUrls() {
	// Choose the folder to serve
	resourceDir := "/assets/"
	router.PathPrefix(resourceDir).
		Handler(http.StripPrefix(resourceDir, http.FileServer(http.Dir("."+resourceDir))))

	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.HandleFunc("/", server.IndexHandler)
	router.HandleFunc("/signup", server.SignupHandler)
	router.HandleFunc("/login", server.LoginHandler)
	router.HandleFunc("/logout", server.LogoutHandler)
	router.HandleFunc("/buyer", buyer.BuyerHandler)
	router.HandleFunc("/seller", seller.SellerHandler)
	router.HandleFunc("/admin", admin.AdminHandler)

	router.HandleFunc("/buyer/findoneitem", buyer.LookForItemHandler)
	router.HandleFunc("/buyer/buyercart", buyer.CartHandler)
	router.HandleFunc("/buyer/checkoutsuccess", buyer.CheckoutSuccessHandler)
	router.HandleFunc("/buyer/profile", buyer.ShowProfile)

	router.HandleFunc("/seller/additem", seller.AddItemHandler)
	router.HandleFunc("/seller/updateitem", seller.UpdateItemHandler)
	router.HandleFunc("/seller/deleteitem", seller.DeleteItemHandler)
	router.HandleFunc("/seller/profile", seller.ShowProfile)
}
