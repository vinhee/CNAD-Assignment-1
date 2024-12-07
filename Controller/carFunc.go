package Controller

import (
	database "CNAD-Assignment-1/Database"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func DisplayCar(w http.ResponseWriter, r *http.Request) {
	nameSess, _ := store.Get(r, "cookieName")
	userName, _ := nameSess.Values["userName"].(string)
	carList, err := database.GetCarDetails()
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}

	tmpl, err := template.ParseFiles("Pages/VehicleManage/vehicleIndex.html", "Pages/UserManage/navbarmember.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "vehicleIndex.html", map[string]interface{}{
		"Cars":     carList,
		"UserName": userName,
	})
	if err != nil {
		log.Println("Error with server: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

var tierNum = map[string]int{ // for member tier
	"Basic":   1,
	"Premium": 2,
	"VIP":     3,
}

func removeDate(dates []string, dateToRemove string) []string {
	for i, date := range dates {
		if date == dateToRemove {
			return append(dates[:i], dates[i+1:]...)
		}
	}
	return dates
}

func BookCar(w http.ResponseWriter, r *http.Request) {
	nameSess, _ := store.Get(r, "cookieName")
	userName, _ := nameSess.Values["userName"].(string)
	emailSess, _ := store.Get(r, "cookieEmail")
	userEmail, _ := emailSess.Values["userEmail"].(string)
	userList, _ := database.GetUserDetail(userEmail)
	var userTier string
	var userBook int
	for _, checkUser := range userList {
		if checkUser.Email == userEmail {
			userTier = checkUser.MemberTier
			userBook = checkUser.Bookings
		}
	}
	log.Print("User bookings: ", userBook)
	stringCarID := r.URL.Query().Get("id")
	if stringCarID == "" {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}
	carID, err := strconv.Atoi(stringCarID)
	if err != nil {
		http.Error(w, "Invalid Car ID", http.StatusBadRequest)
		return
	}
	log.Print("Car ID:", carID)
	car, err := database.GetSpecificCar(carID)
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}

	bookList, _ := database.GetCarBook(carID)
	log.Print("book details: ", bookList)

	var blockedDates []string
	for _, booking := range bookList {
		if booking.Status != "Cancelled" {
			startDate := booking.StartDate
			endDate := booking.EndDate
			for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
				blockedDates = append(blockedDates, d.Format("2006-01-02"))
			}
		} else {
			startDate := booking.StartDate
			endDate := booking.EndDate
			for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
				blockedDates = removeDate(blockedDates, d.Format("2006-01-02"))
			}
		}
	}

	blockedDatesJSON, err := json.Marshal(blockedDates)
	if err != nil {
		log.Println("Error marshaling blocked dates:", err)
		return
	}
	log.Print("blocked dates JSON: ", string(blockedDatesJSON))

	today := time.Now().Format("2006-01-02")

	userTierNum := tierNum[userTier]
	carTierNum := tierNum[car.MemberTier]

	if userTierNum >= carTierNum {
		tmpl, err := template.ParseFiles("Pages/VehicleManage/bookcarpage.html", "Pages/UserManage/navbarmember.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "bookcarpage.html", map[string]interface{}{
			"Cars":         car,
			"UserName":     userName,
			"Today":        today,
			"BlockedDates": string(blockedDatesJSON),
		})
		if err != nil {
			log.Println("Error with server: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else if userTierNum < carTierNum {
		errMsg := "Your membership tier does not meet the minimum membership tier"
		tmpl, err := template.ParseFiles("Pages/VehicleManage/bookcarpage.html", "Pages/UserManage/navbarmember.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "bookcarpage.html", map[string]interface{}{
			"Cars":     car,
			"UserName": userName,
			"Today":    today,
			"ErrMsg":   errMsg,
		})
		if err != nil {
			log.Println("Error with server: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if userBook == 0 {
		errBookMsg := "You have exceeded your booking limits!"
		tmpl, err := template.ParseFiles("Pages/VehicleManage/bookcarpage.html", "Pages/UserManage/navbarmember.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "bookcarpage.html", map[string]interface{}{
			"Cars":       car,
			"UserName":   userName,
			"Today":      today,
			"ErrBookMsg": errBookMsg,
		})
		if err != nil {
			log.Println("Error with server: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func roundNum(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	IdSess, _ := store.Get(r, "cookieID")
	userID := IdSess.Values["userID"]
	emailSess, _ := store.Get(r, "cookieEmail")
	userEmail := emailSess.Values["userEmail"]
	log.Print("UserID: ", userID)
	carName := r.FormValue("carName")

	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")
	pickupTime := r.FormValue("pickupTime")
	dropoffTime := r.FormValue("dropoffTime")

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

	var priceHour int
	var imagelink string
	carList, err := database.GetCarDetails()
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}
	for _, findCar := range carList {
		if findCar.Name == carName {
			priceHour = findCar.PriceHour
			imagelink = findCar.ImageLink
			break
		}
	}

	startDateTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", startDate, pickupTime))
	if err != nil {
		log.Println("Error parsing start date and time:", err)
		return
	}

	endDateTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", endDate, dropoffTime))
	if err != nil {
		log.Println("Error parsing end date and time:", err)
		return
	}

	totalHours := endDateTime.Sub(startDateTime).Hours()
	totalCost := totalHours * float64(priceHour)

	totalCostString := fmt.Sprintf("%.2f Hours x $%d (Standard Price per Hour) =  $%.2f ", totalHours, priceHour, totalCost)

	formattedStartDateTime := startDateTime.Format("2006-01-02")
	formattedEndDateTime := endDateTime.Format("2006-01-02")

	carIDStr := r.FormValue("carID")
	if carIDStr == "" {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}
	carID, err := strconv.Atoi(carIDStr)
	if err != nil {
		http.Error(w, "Invalid Car ID", http.StatusBadRequest)
		return
	}

	var booking database.CarsBooking
	booking.UserID = userID.(int)
	booking.CarName = carName
	booking.CarID = carID
	booking.StartDate = startDateTime
	booking.EndDate = endDateTime
	booking.TotalHours = roundNum(totalHours, 2) // rounds number to 2 d.p.
	booking.TotalCost = roundNum(totalCost, 2)

	database.AddBooking(booking) // add to booking database
	database.UpdateUserBook(userEmail.(string))

	tmpl, err := template.ParseFiles("Pages/VehicleManage/confirmbookpage.html", "Pages/UserManage/navbarmember.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "confirmbookpage.html", map[string]interface{}{
		"UserName":        userName,
		"StartDateTime":   formattedStartDateTime,
		"EndDateTime":     formattedEndDateTime,
		"CarName":         carName,
		"ImageLink":       imagelink,
		"TotalCostString": totalCostString,
	})
	if err != nil {
		log.Println("Error with server: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func EditBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
		startDate := startDateTime.Format("2006-01-02")
		startTime := startDateTime.Format("15:04")
		endDate := endDateTime.Format("2006-01-02")
		endTime := endDateTime.Format("15:04")

		userID := book.UserID
		carID := book.CarID
		log.Print("Car booking for edit:", carID, userID, bookingID, startDateTime, endDateTime, startDate, startTime, endDate, endTime)

		car, _ := database.GetSpecificCar(carID)
		carImage := car.ImageLink
		carName := car.Name

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

		tmpl, err := template.ParseFiles("Pages/VehicleManage/editbookingpage.html", "Pages/UserManage/navbarmember.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "editbookingpage.html", map[string]interface{}{
			"UserName":  userName,
			"StartDate": startDate,
			"StartTime": startTime,
			"EndDate":   endDate,
			"EndTime":   endTime,
			"CarName":   carName,
			"ImageLink": carImage,
			"CarID":     carID,
			"BookingID": bookingID,
		})
		if err != nil {
			log.Println("Error with server nooooo: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodPost {
		bookingIDstr := r.FormValue("bookingID")
		log.Println("Received bookingID:", bookingIDstr)
		carIDstr := r.FormValue("carID")
		log.Print("CarID for method post: ", carIDstr)
		carName := r.FormValue("carName")
		startDate := r.FormValue("startDate")
		endDate := r.FormValue("endDate")
		startTime := r.FormValue("pickupTime")
		endTime := r.FormValue("dropoffTime")

		bookingID, err := strconv.Atoi(bookingIDstr)
		if err != nil {
			http.Error(w, "Invalid booking ID", http.StatusBadRequest)
			return
		}

		carID, err := strconv.Atoi(carIDstr)
		if err != nil {
			log.Println("Error converting carID to integer:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		book, _ := database.GetBookingByID(bookingID)
		userID := book.UserID

		log.Print("Car booking for update:", carID, userID, bookingID, startDate, startTime, endDate, endTime)

		startDateTimeStr := startDate + " " + startTime
		endDateTimeStr := endDate + " " + endTime

		startDateTime, err := time.Parse("2006-01-02 15:04", startDateTimeStr)
		if err != nil {
			log.Println("Error parsing start date and time:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		endDateTime, err := time.Parse("2006-01-02 15:04", endDateTimeStr)
		if err != nil {
			log.Println("Error parsing end date and time:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		car, _ := database.GetSpecificCar(carID)
		priceHour := car.PriceHour

		totalHours := endDateTime.Sub(startDateTime).Hours()
		totalCost := totalHours * float64(priceHour)

		err = database.UpdateCarBook(userID, carName, carID, startDateTime, endDateTime, totalHours, totalCost, "Booked", bookingID)
		if err != nil {
			log.Println("Error updating booking:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}
