package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-gorp/gorp"

	_ "github.com/mattn/go-sqlite3"
)

var db *gorp.DbMap

func main() {
	db = initDB("./goto.db")
	defer db.Db.Close()

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
