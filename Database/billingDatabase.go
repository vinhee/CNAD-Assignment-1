package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Billing struct {
	Id        int       `json:"id"`
	UserID    int       `json:"userid"`
	CarID     int       `json:"carid"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
	PriceHour int       `json:"priceHour"`
	TotalCost float64   `json:"totalcost"`
	Status    string    `json:"status"`
}

var billdb *sql.DB

func InsertBill(bill Billing) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	defer db.Close()
	insertQuery := "INSERT INTO Billing (UserID, CarID, StartDate, EndDate, PriceHour, TotalCost, Status) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, bill.UserID, bill.CarID, bill.StartDate, bill.EndDate, bill.PriceHour, bill.TotalCost, "Paid")
	if err != nil {
		log.Println("Database insert error:", err)
		return err
	}
	return nil
}

func Bill() {

}
