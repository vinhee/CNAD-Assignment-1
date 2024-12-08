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
	BookingID int       `json:"bookid"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
	PriceHour int       `json:"priceHour"`
	TotalCost float64   `json:"totalcost"`
	UserCard  string    `json:"usercard"`
	Status    string    `json:"status"`
}

var billdb *sql.DB

func InsertBill(bill Billing) (billID int64, err error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return 0, err
	}
	defer db.Close()

	insertQuery := "INSERT INTO Billing (UserID, CarID, BookingID, StartDate, EndDate, PriceHour, TotalCost, UserCard, Status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(insertQuery, bill.UserID, bill.CarID, bill.BookingID, bill.StartDate, bill.EndDate, bill.PriceHour, bill.TotalCost, "", bill.Status)
	if err != nil {
		log.Println("Database insert error:", err)
		return 0, err
	}

	billID, err = result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return 0, err
	}

	return billID, nil
}

func GetBill(billID int) (Billing, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return Billing{}, err
	}
	defer db.Close()

	query := "SELECT * FROM Billing WHERE ID = ?"
	row := db.QueryRow(query, billID)

	var bill Billing
	var startDateBytes, endDateBytes []byte

	if err := row.Scan(&bill.Id, &bill.UserID, &bill.CarID, &bill.BookingID, &startDateBytes, &endDateBytes, &bill.PriceHour, &bill.TotalCost, &bill.UserCard, &bill.Status); err != nil {
		log.Println("Row scan error:", err)
		return Billing{}, err
	}

	startDateTime, err := time.Parse("2006-01-02 15:04:05", string(startDateBytes))
	if err != nil {
		log.Println("Error parsing StartDate:", err)
		return Billing{}, err
	}
	bill.StartDate = startDateTime

	endDateTime, err := time.Parse("2006-01-02 15:04:05", string(endDateBytes))
	if err != nil {
		log.Println("Error parsing EndDate:", err)
		return Billing{}, err
	}
	bill.EndDate = endDateTime

	return bill, nil
}

func UpdateBillCard(billID int, userCard string) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE Billing SET UserCard = ?, Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, userCard, "Paid", billID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}
