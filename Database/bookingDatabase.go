package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type CarsBooking struct {
	Id         int       `json:"id"`
	UserID     int       `json:"userid"`
	CarName    string    `json:"carname"`
	StartDate  time.Time `json:"startdate"`
	EndDate    time.Time `json:"enddate"`
	TotalHours float64   `json:"totalhours"`
	TotalCost  float64   `json:"totalcost"`
	Status     string    `json:"status"`
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
	insertQuery := "INSERT INTO CarsBooking (UserID, CarName, StartDate, EndDate, TotalHours, TotalCost, Status) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, booking.UserID, booking.CarName, booking.StartDate, booking.EndDate, booking.TotalHours, booking.TotalCost, "Booked")
	if err != nil {
		log.Println("Database insert error:", err)
		return err
	}
	return nil
}

func UpdatePaid(bookingID int) error {
	db, err := GetBookingDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE Users SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "Paid", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func UpdateComplete(bookingID int) error {
	db, err := GetBookingDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE Users SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "Completed", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func UpdateCancelled(bookingID int) error {
	db, err := GetBookingDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE Users SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "Cancelled", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func GetUserBook(userID int) ([]CarsBooking, error) {
	db, err := GetBookingDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM CarsBooking WHERE UserID = ?"
	results, err := db.Query(query, userID)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	bookList := []CarsBooking{}
	for results.Next() {
		var book CarsBooking
		var startDate, endDate []byte
		if err := results.Scan(&book.Id, &book.UserID, &book.CarName, &startDate, &endDate, &book.TotalHours, &book.TotalCost, &book.Status); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		book.StartDate, err = parseDateTime(startDate)
		if err != nil {
			log.Println("Error parsing StartDate:", err)
			return nil, err
		}

		book.EndDate, err = parseDateTime(endDate)
		if err != nil {
			log.Println("Error parsing EndDate:", err)
			return nil, err
		}
		bookList = append(bookList, book)
	}
	for _, book := range bookList {
		fmt.Printf("ID: %d, UserID: %d, CarName: %s, StartDate: %s, EndDate: %s, TotalHours: %.2f, TotalCost: %.2f, Status: %s\n",
			book.Id, book.UserID, book.CarName, book.StartDate, book.EndDate, book.TotalHours, book.TotalCost, book.Status)
	}
	return bookList, nil
}

func parseDateTime(data []byte) (time.Time, error) {
	dateStr := string(data)
	return time.Parse("2006-01-02 15:04:05", dateStr)

}
