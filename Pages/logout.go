package Pages

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookieName")
	session.Values["cookieName"] = nil
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)
}
