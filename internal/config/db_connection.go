package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func SetupDB() {
	var err error

	dbString := "root:@tcp(localhost:3306)/my_budget"

	Db, err = sql.Open("mysql", dbString)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}

	// Check if the MySQL database connection is alive
	err = Db.Ping()
	if err != nil {
		fmt.Println("Error checking database connection:", err)
		return
	}
	fmt.Println("Successfully connected to the database!")
}

// Function to close the MySQL database connection
func CloseDB() {
	if Db != nil {
		Db.Close()
	}
}
