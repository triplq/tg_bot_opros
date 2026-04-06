package main

import (
	"database/sql"
	"fmt"
	"log"
	"parser/config"
)

func connect2db() *sql.DB {
	user := config.Config("POSTGRES_USER")
	pwd := config.Config("POSTGRES_PASSWORD")
	dbname := config.Config("POSTGRES_DB")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pwd, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := connect2db()
}
