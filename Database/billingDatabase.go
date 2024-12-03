package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Billing struct {
	Id         int    `json:"id"`
	UserID     int    `json:"userid"`
	CarName    string `json:"carname"`
	StartDate  string `json:"startdate"`
	EndDate    string `json:"enddate"`
	TotalHours int    `json:"totalhours"`
	TotalCost  int    `json:"totalcost"`
}

var billdb *sql.DB

func GetBillingDB() (*sql.DB, error) {
	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	return db, nil
}
