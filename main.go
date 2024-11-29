package main

import (
	"CNAD-Assignment-1/Controller"
	"net/http"
)

func main() {
	http.HandleFunc("/login", Controller.Loginpage)
	http.HandleFunc("/homepage", Controller.HomePage)
	http.HandleFunc("/register", Controller.Registerpage)
	http.HandleFunc("/homemember", Controller.HomeMember)
	http.HandleFunc("/logout", Controller.Logout)

	http.ListenAndServe(":5000", nil)
}
