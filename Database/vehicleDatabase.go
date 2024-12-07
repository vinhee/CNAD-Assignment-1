package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Cars struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageLink   string `json:"imagelink"`
	PriceHour   int    `json:"pricehour"`
	MemberTier  string `json:"membertier"`
}

var cardb *sql.DB

func GetCarDetails() ([]Cars, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM Cars"
	results, err := db.Query(query)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	carList := []Cars{}
	for results.Next() {
		var car Cars
		if err := results.Scan(&car.Id, &car.Name, &car.Description, &car.ImageLink, &car.PriceHour, &car.MemberTier); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		carList = append(carList, car)
	}
	return carList, nil
}

func GetSpecificCar(carID int) (Cars, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return Cars{}, err
	}
	defer db.Close()
	query := "SELECT * FROM Cars WHERE ID = ?"
	row := db.QueryRow(query, carID)

	var car Cars
	if err := row.Scan(&car.Id, &car.Name, &car.Description, &car.ImageLink, &car.PriceHour, &car.MemberTier); err != nil {
		log.Println("Row scan error:", err)
		return Cars{}, err
	}

	return car, nil
}
