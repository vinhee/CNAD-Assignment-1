package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type CarsBooking struct {
	Id         int     `json:"id"`
	UserID     int     `json:"userid"`
	CarName    string  `json:"carname"`
	StartDate  string  `json:"startdate"`
	EndDate    string  `json:"enddate"`
	TotalHours float64 `json:"totalhours"`
	TotalCost  float64 `json:"totalcost"`
}

var bookdb *sql.DB

func GetBookingDB() (*sql.DB, error) {
	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	return db, nil
}

func AddBooking(booking CarsBooking) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	defer db.Close()
	insertQuery := "INSERT INTO CarsBooking (UserID, CarName, StartDate, EndDate, TotalHours, TotalCost) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, booking.UserID, booking.CarName, booking.StartDate, booking.EndDate, booking.TotalHours, booking.TotalCost)
	if err != nil {
		log.Println("Database insert error:", err)
		return err
	}
	return nil
}
