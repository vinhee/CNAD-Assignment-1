package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type CarsBooking struct {
	Id         int       `json:"id"`
	UserID     int       `json:"userid"`
	CarName    string    `json:"carname"`
	CarID      int       `json:"carid"`
	StartDate  time.Time `json:"startdate"`
	EndDate    time.Time `json:"enddate"`
	TotalHours float64   `json:"totalhours"`
	TotalCost  float64   `json:"totalcost"`
	Status     string    `json:"status"`
}

var bookdb *sql.DB

func AddBooking(booking CarsBooking) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	defer db.Close()
	insertQuery := "INSERT INTO CarsBooking (UserID, CarName, CarID, StartDate, EndDate, TotalHours, TotalCost, Status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, booking.UserID, booking.CarName, booking.CarID, booking.StartDate, booking.EndDate, booking.TotalHours, booking.TotalCost, "Booked")
	if err != nil {
		log.Println("Database insert error:", err)
		return err
	}
	return nil
}

func GetBookingByID(bookingID int) (CarsBooking, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return CarsBooking{}, err
	}
	defer db.Close()

	query := "SELECT * FROM CarsBooking WHERE ID = ?"
	row := db.QueryRow(query, bookingID)

	var book CarsBooking
	var startDateBytes, endDateBytes []byte

	if err := row.Scan(&book.Id, &book.UserID, &book.CarName, &book.CarID, &startDateBytes, &endDateBytes, &book.TotalHours, &book.TotalCost, &book.Status); err != nil {
		log.Println("Row scan error:", err)
		return CarsBooking{}, err
	}

	startDateTime, err := time.Parse("2006-01-02 15:04:05", string(startDateBytes))
	if err != nil {
		log.Println("Error parsing StartDate:", err)
		return CarsBooking{}, err
	}
	book.StartDate = startDateTime

	endDateTime, err := time.Parse("2006-01-02 15:04:05", string(endDateBytes))
	if err != nil {
		log.Println("Error parsing EndDate:", err)
		return CarsBooking{}, err
	}
	book.EndDate = endDateTime

	return book, nil
}

func UpdatePaid(bookingID int) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE CarsBooking SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "Paid", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func UpdateInProgress(bookingID int) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE CarsBooking SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "In Progress", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func UpdateComplete(bookingID int) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE CarsBooking SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "Completed", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func UpdateCancelled(bookingID int) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE CarsBooking SET Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, "Cancelled", bookingID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}

func GetUserBook(userID int) ([]CarsBooking, error) {
	db, err := GetDB()
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
		if err := results.Scan(&book.Id, &book.UserID, &book.CarName, &book.CarID, &startDate, &endDate, &book.TotalHours, &book.TotalCost, &book.Status); err != nil {
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
		fmt.Printf("ID: %d, UserID: %d, CarName: %s, CarID: %d, StartDate: %s, EndDate: %s, TotalHours: %.2f, TotalCost: %.2f, Status: %s\n",
			book.Id, book.UserID, book.CarName, book.CarID, book.StartDate, book.EndDate, book.TotalHours, book.TotalCost, book.Status)
	}
	return bookList, nil
}

func parseDateTime(data []byte) (time.Time, error) {
	dateStr := string(data)
	return time.Parse("2006-01-02 15:04:05", dateStr)

}

func GetCarBook(carID int) ([]CarsBooking, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM CarsBooking WHERE CarID = ?"
	results, err := db.Query(query, carID)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	bookList := []CarsBooking{}
	for results.Next() {
		var book CarsBooking
		var startDate, endDate []byte
		if err := results.Scan(&book.Id, &book.UserID, &book.CarName, &book.CarID, &startDate, &endDate, &book.TotalHours, &book.TotalCost, &book.Status); err != nil {
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

	return bookList, nil
}

func UpdateCarBook(userID int, carName string, carID int, startDateTime time.Time, endDateTime time.Time, totalHours float64, totalCost float64, status string, bookID int) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE CarsBooking SET UserID = ?, CarName = ?, CarID = ?, StartDate = ?, EndDate = ?, TotalHours = ?, TotalCost = ?, Status = ? WHERE ID = ?"
	_, err = db.Exec(updateQuery, userID, carName, carID, startDateTime, endDateTime, totalHours, totalCost, status, bookID)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}
