package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func insertFileNamesIntoDB(files []string) error {
	// Open the SQLite database.
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Create a new table.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS strings (id INTEGER PRIMARY KEY, value TEXT)")
	if err != nil {
		return err
	}

	// Insert the strings into the table.
	for _, value := range files {
		_, err = db.Exec("INSERT INTO strings (value) VALUES (?)", value)
		if err != nil {
			return err
		}
	}
	return nil
}

func countEntry() error {
	// Open the SQLite database.
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the number of values in the table.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM strings").Scan(&count)
	if err != nil {
		return err
	}
	fmt.Printf("There are %d values in the table.\n", count)

	return nil
}
