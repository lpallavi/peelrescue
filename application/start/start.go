package start

import (
	"log"
	"net/http"

	config "projectGoLive/application/config"

	"github.com/gorilla/mux"
)

var (
	router = mux.NewRouter()
)

func StartApplication() {

	mapUrls()
	log.Println(" Listening on port ", config.PortNum)
	log.Fatal(http.ListenAndServeTLS(config.PortNum, config.CertPath+"cert.pem", config.CertPath+"key.pem", router))
}
