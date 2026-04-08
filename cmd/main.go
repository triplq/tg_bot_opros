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

	tg, err := clients.Init()
	if err != nil {
		log.Fatal(err)
	}

	var url = []string{"https://t.me/s/oprosrussiaa", "https://t.me/s/oprosyvsem", "https://t.me/s/servizoria_rf"}

	app := &application.App{
		Model: &database.Model{DB: db},
		Tg:    tg,
	}

	for _, u := range url {
		err = parsing.Parse(u, app)
		if err != nil {
			log.Fatal(err)
		}
	}
}
