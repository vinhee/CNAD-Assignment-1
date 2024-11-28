package main

import (
	"CNAD-Assignment-1/Pages"
	"net/http"
)

func main() {
	http.HandleFunc("/login", Pages.Loginpage)
	http.HandleFunc("/homepage", Pages.Homepage)
	http.HandleFunc("/register", Pages.Registerpage)
	http.HandleFunc("/homemember", Pages.HomeMember)
	http.HandleFunc("/logout", Pages.Logout)

	http.ListenAndServe(":5000", nil)
}
