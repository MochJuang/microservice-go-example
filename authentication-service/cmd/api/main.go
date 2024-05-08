package main

import (
	"authentication-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect db
	conn := connectTODB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// init config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// initial server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectTODB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready ...")
			counts++

			if counts > 10 {
				log.Println(err.Error())
				return nil
			}
		}

		if err == nil {
			log.Println("Connected to Postgres!")
			return connection
		}

		log.Println("Backing off for seconds... ")
		time.Sleep(2 * time.Second)

		continue
	}

}
