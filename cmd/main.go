package main

import (
	"database/sql"
	"fmt"
	"log"
	"parser/application"
	"parser/clients"
	"parser/config"
	"parser/database"
	"parser/parsing"

	_ "github.com/lib/pq"
)

func connect2db() *sql.DB {
	user := config.Config("POSTGRES_USER")
	pwd := config.Config("POSTGRES_PASSWORD")
	dbname := config.Config("POSTGRES_DB")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=localhost port=5432", user, pwd, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := connect2db()
	defer db.Close()

	_, err := clients.Init()
	if err != nil {
		log.Fatal(err)
	}

	url := "https://t.me/s/oprosrussiaa"

	app := &application.App{
		Model: &database.Model{DB: db},
		Tg:    tg,
	}

	parsing.Parse(url, app)
}
