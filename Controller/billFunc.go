package Controller

import (
	database "CNAD-Assignment-1/Database"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tierDisc = map[string]float64{ // for member tier
	"Basic":   0.95,
	"Premium": 0.90,
	"VIP":     0.85,
}

func DisplayBill(w http.ResponseWriter, r *http.Request) {
	bookingIDstr := r.URL.Query().Get("bookingID")
	bookingID, err := strconv.Atoi(bookingIDstr)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	book, err := database.GetBookingByID(bookingID)
	if err != nil {
		log.Println("Error retrieving booking:", err)

		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	startDateTime := book.StartDate
	endDateTime := book.EndDate
	totalHours := book.TotalHours

	formattedStartDateTime := startDateTime.Format("2006-01-02")
	formattedEndDateTime := endDateTime.Format("2006-01-02")

	userID := book.UserID
	carID := book.CarID
	log.Print("Car booking for payment:", carID, userID, bookingID, startDateTime, endDateTime)

	car, _ := database.GetSpecificCar(carID)
	carImage := car.ImageLink
	carName := car.Name
	priceHour := car.PriceHour

	var userName string
	var userTier string
	userList, err := database.GetLoginUser()
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}
	for _, findUser := range userList {
		if findUser.Id == userID {
			userName = findUser.Name
			userTier = findUser.MemberTier
			break
		}
	}

	memberDisc := tierDisc[userTier]
	log.Print("Member Discount: ", memberDisc)

	cost := totalHours * float64(priceHour)
	totalCost := cost * memberDisc

	costString := fmt.Sprintf("%.2f Hours x $%d (Price per Hour) =  $%.2f ", totalHours, priceHour, cost)

	totalCoststr := fmt.Sprintf("%s x %.2f  =  $%.2f ", costString, memberDisc, totalCost)

	memberDiscStr := fmt.Sprintf("%s = %.2f", userTier, memberDisc)

	var bill database.Billing
	bill.UserID = userID
	bill.CarID = carID
	bill.BookingID = bookingID
	bill.StartDate = startDateTime
	bill.EndDate = endDateTime
	bill.PriceHour = priceHour // rounds number to 2 d.p.
	bill.TotalCost = roundNum(totalCost, 2)
	bill.Status = "Unpaid"
	bill.CarName = carName

	billID, err := database.InsertBill(bill) // insert new bill'
	if err != nil {
		log.Print("Error adding new bill: ", err)
	}

	tmpl, err := template.ParseFiles("Pages/BillManage/paymentpage.html", "Pages/UserManage/navbarmember.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "paymentpage.html", map[string]interface{}{
		"UserName":      userName,
		"StartDateTime": formattedStartDateTime,
		"EndDateTime":   formattedEndDateTime,
		"CarName":       carName,
		"ImageLink":     carImage,
		"CarID":         carID,
		"BookingID":     bookingID,
		"CostString":    costString,
		"MemberDisc":    memberDiscStr,
		"TotalCost":     totalCoststr,
		"BillID":        billID,
	})
	if err != nil {
		log.Println("Error with server in Display Bill: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	billIDstr := r.URL.Query().Get("billID")
	billID, err := strconv.Atoi(billIDstr)
	if err != nil {
		http.Error(w, "Invalid billing ID", http.StatusBadRequest)
		return
	}

	bill, err := database.GetBill(billID)
	carID := bill.CarID
	userID := bill.UserID
	startDateTime := bill.StartDate
	endDateTime := bill.EndDate
	bookingID := bill.BookingID
	carName := bill.CarName
	log.Print("Booking ID for confirm: ", bookingID)

	formattedStartDateTime := startDateTime.Format("2006-01-02")
	formattedEndDateTime := endDateTime.Format("2006-01-02")

	totalCost := bill.TotalCost

	car, _ := database.GetSpecificCar(carID)
	carImage := car.ImageLink

	var userName string
	userList, err := database.GetLoginUser()
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}
	for _, findUser := range userList {
		if findUser.Id == userID {
			userName = findUser.Name
			break
		}
	}

	tmpl, err := template.ParseFiles("Pages/BillManage/makepaymentpage.html", "Pages/UserManage/navbarmember.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "makepaymentpage.html", map[string]interface{}{
		"UserName":  userName,
		"StartDate": formattedStartDateTime,
		"EndDate":   formattedEndDateTime,
		"CarName":   carName,
		"ImageLink": carImage,
		"TotalCost": totalCost,
		"BillID":    billID,
		"BookingID": bookingID,
	})
	if err != nil {
		log.Println("Error with server in Display Bill: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func UpdatePaymentCard(w http.ResponseWriter, r *http.Request) {
	billIDstr := r.FormValue("billID")
	billID, err := strconv.Atoi(billIDstr)
	if err != nil {
		http.Error(w, "Invalid billing ID", http.StatusBadRequest)
		return
	}
	bookIDstr := r.FormValue("bookingID")
	bookID, err := strconv.Atoi(bookIDstr)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		log.Print("Booking ID problem: ", bookID)
		return
	}
	userCard := r.FormValue("userCard")

	database.UpdateBillCard(billID, userCard)
	database.UpdatePaid(bookID)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
