//Admin functions
package admin

import (
	"net/http"
	"projectGoLive/application/config"
	"projectGoLive/application/server"
)

//---------------------------------------------------------------------------
// Functions for admin
//---------------------------------------------------------------------------
// This method is used to perform all admin functions
func AdminHandler(w http.ResponseWriter, req *http.Request) {
	if !server.ActiveSession(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	config.TPL.ExecuteTemplate(w, "admin.gohtml", nil)
}
