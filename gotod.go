package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	initDB("./goto.db")
	defer db.Close()

	router := NewRouter()
	server := http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	fmt.Println("Listening on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initDB(file string) {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}
}
