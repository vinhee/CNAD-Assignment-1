package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Billing struct {
	Id        int       `json:"id"`
	UserID    int       `json:"userid"`
	CarID     int       `json:"carid"`
	BookingID int       `json:"bookid"`
	CarName   string    `json:"carname"`
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

	insertQuery := "INSERT INTO Billing (UserID, CarID, BookingID, CarName, StartDate, EndDate, PriceHour, TotalCost, UserCard, Status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(insertQuery, bill.UserID, bill.CarID, bill.BookingID, bill.CarName, bill.StartDate, bill.EndDate, bill.PriceHour, bill.TotalCost, "", bill.Status)
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

	if err := row.Scan(&bill.Id, &bill.UserID, &bill.CarID, &bill.BookingID, &bill.CarName, &startDateBytes, &endDateBytes, &bill.PriceHour, &bill.TotalCost, &bill.UserCard, &bill.Status); err != nil {
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

func GetBillByUser(userID int) ([]Billing, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM Billing WHERE UserID = ?"
	results, err := db.Query(query, userID)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	billList := []Billing{}
	for results.Next() {
		var bill Billing
		var startDate, endDate []byte
		if err := results.Scan(&bill.Id, &bill.UserID, &bill.CarID, &bill.BookingID, &bill.CarName, &startDate, &endDate, &bill.PriceHour, &bill.TotalCost, &bill.UserCard, &bill.Status); err != nil {
			log.Println("Row scan error for billing db:", err)
			return nil, err
		}
		bill.StartDate, err = parseDateTime(startDate)
		if err != nil {
			log.Println("Error parsing StartDate:", err)
			return nil, err
		}

		bill.EndDate, err = parseDateTime(endDate)
		if err != nil {
			log.Println("Error parsing EndDate:", err)
			return nil, err
		}
		billList = append(billList, bill)
	}

	for _, bill := range billList {
		fmt.Printf("ID: %d, UserID: %d, CarID: %d, BookID: %d, CarName: %s, StartDate: %s, EndDate: %s, PriceHour: %d, TotalCost: %.2f, UserCard: %s, Status: %s\n",
			bill.Id, bill.UserID, bill.CarID, bill.BookingID, bill.CarName, bill.StartDate, bill.EndDate, bill.PriceHour, bill.TotalCost, bill.UserCard, bill.Status)
	}

	return billList, nil
}
