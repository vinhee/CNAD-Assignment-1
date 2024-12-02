package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Cars struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageLink   string `json:"imagelink"`
	PriceHour   int    `json:"pricehour"`
	Quantity    int    `json:"quantity"`
	MemberTier  string `json:"membertier"`
}

var cardb *sql.DB

func GetCarDB() (*sql.DB, error) {
	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	cardb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	return cardb, nil
}

func GetCarDetails() ([]Cars, error) {
	db, err := GetCarDB()
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
		if err := results.Scan(&car.Id, &car.Name, &car.Description, &car.ImageLink, &car.PriceHour, &car.Quantity, &car.MemberTier); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		carList = append(carList, car)
	}
	return carList, nil
}

func GetSpecificCar(carID int) ([]Cars, error) {
	db, err := GetCarDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM Cars WHERE ID = ?"
	results, err := db.Query(query, carID)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	carList := []Cars{}
	for results.Next() {
		var car Cars
		if err := results.Scan(&car.Id, &car.Name, &car.ImageLink, &car.PriceHour, &car.Quantity, &car.MemberTier); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		carList = append(carList, car)
	}
	return carList, nil
}
