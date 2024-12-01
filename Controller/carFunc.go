package Controller

import (
	database "CNAD-Assignment-1/Database"
	"html/template"
	"log"
	"net/http"
)

func DisplayCar(w http.ResponseWriter, r *http.Request) {
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
		"Cars": carList,
	})
	if err != nil {
		log.Println("Error with server: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
