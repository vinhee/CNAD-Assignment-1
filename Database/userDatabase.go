package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	MemberTier string `json:"memberTier"`
}

var db *sql.DB

func GetDB() (*sql.DB, error) {
	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	return db, nil
}

func GetLoginUser() ([]User, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM Users"
	results, err := db.Query(query)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	userList := []User{}
	for results.Next() {
		var user User
		if err := results.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.MemberTier); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func InsertNewUser(user User) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	defer db.Close()
	insertQuery := "INSERT INTO Users (Name, Email, Password, MemberTier) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, user.Name, user.Email, user.Password, "Basic")
	if err != nil {
		log.Println("Database insert error:", err)
		return err
	}
	return nil
}

func GetUserDetail(userEmail string) ([]User, error) {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return nil, err
	}
	defer db.Close()
	query := "SELECT * FROM Users WHERE Email = ?"
	results, err := db.Query(query, userEmail)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer results.Close()

	userList := []User{}
	for results.Next() {
		var user User
		if err := results.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.MemberTier); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func UpdateUser(userName string, userEmail string, userPassword string, userTier string) error {
	db, err := GetDB()
	if err != nil {
		log.Println("Unable to connect to function:", err)
		return err
	}
	updateQuery := "UPDATE Users SET Name = ?, Email = ?, Password = ?, MemberTier = ? WHERE Email = ?"
	_, err = db.Exec(updateQuery, userName, userEmail, userPassword, userTier, userEmail)
	if err != nil {
		log.Println("Database update error:", err)
		return err
	}
	return nil
}
