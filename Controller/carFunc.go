package Controller

import (
	database "CNAD-Assignment-1/Database"
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

func BookCar(w http.ResponseWriter, r *http.Request) {
	nameSess, _ := store.Get(r, "cookieName")
	userName, _ := nameSess.Values["userName"].(string)
	emailSess, _ := store.Get(r, "cookieEmail")
	userEmail, _ := emailSess.Values["userEmail"].(string)
	userList, _ := database.GetUserDetail(userEmail)
	var userTier string
	for _, checkUser := range userList {
		if checkUser.Email == userEmail {
			userTier = checkUser.MemberTier
		}
	}
	stringcarID := r.URL.Query().Get("id")
	if stringcarID == "" {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}
	carID, _ := strconv.Atoi(stringcarID)

	car, err := database.GetSpecificCar(carID)
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}

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
			"Cars":     car,
			"UserName": userName,
			"Today":    today,
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
	var carID int
	carList, err := database.GetCarDetails()
	if err != nil {
		log.Println("Retrieve Data Error:", err)
		return
	}
	for _, findCar := range carList {
		if findCar.Name == carName {
			priceHour = findCar.PriceHour
			imagelink = findCar.ImageLink
			carID = findCar.Id
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

	var booking database.CarsBooking
	booking.UserID = userID.(int)
	booking.CarName = carName
	booking.StartDate = startDateTime
	booking.EndDate = endDateTime
	booking.TotalHours = roundNum(totalHours, 2) // rounds number to 2 d.p.
	booking.TotalCost = roundNum(totalCost, 2)

	database.AddBooking(booking)   // add to booking database
	database.UpdateQuantity(carID) // update quantity of car to 0
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
