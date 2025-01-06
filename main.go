package main

import (
	"CNAD-Assignment-1/Controller"
	"net/http"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// User Management
	http.HandleFunc("/login", Controller.Loginpage)
	http.HandleFunc("/register", Controller.Registerpage)
	http.HandleFunc("/homemember", Controller.HomeMember)
	http.HandleFunc("/logout", Controller.Logout)
	http.HandleFunc("/profile", Controller.ProfilePage)
	http.HandleFunc("/editprofile", Controller.EditProfile)

	// Vehicle Reservation
	http.HandleFunc("/displaycar", Controller.DisplayCar)
	http.HandleFunc("/bookcar", Controller.BookCar)
	http.HandleFunc("/confirmbooking", Controller.ConfirmBooking)
	http.HandleFunc("/cancelbooking", Controller.CancelBooking)
	http.HandleFunc("/editbooking", Controller.EditBooking)

	// Billing
	http.HandleFunc("/payment", Controller.DisplayBill)
	http.HandleFunc("/makepayment", Controller.ConfirmPayment)
	http.HandleFunc("/updatepayment", Controller.UpdatePaymentCard)

	http.ListenAndServe(":5000", nil)
}
