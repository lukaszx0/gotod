package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-gorp/gorp"

	_ "github.com/mattn/go-sqlite3"
)

var db *gorp.DbMap

func main() {
	dbfile := flag.String("file", "./goto.db", "Database file")
	port := flag.Int("port", 8080, "Port")
	flag.Parse()

	db = initDB(*dbfile)
	defer closeDB()

	router := NewRouter()
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	fmt.Printf("Listening on :%d\n", *port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
