package main

import (
	"CNAD-Assignment-1/Controller"
	"net/http"
)

func main() {
	http.HandleFunc("/login", Controller.Loginpage)
	http.HandleFunc("/register", Controller.Registerpage)
	http.HandleFunc("/homemember", Controller.HomeMember)
	http.HandleFunc("/logout", Controller.Logout)
	http.HandleFunc("/profile", Controller.ProfilePage)
	http.HandleFunc("/editprofile", Controller.EditProfile)
	http.HandleFunc("/displaycar", Controller.DisplayCar)
	http.HandleFunc("/bookcar", Controller.BookCar)

	http.ListenAndServe(":5000", nil)
}
