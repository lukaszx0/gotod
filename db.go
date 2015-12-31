package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
)

// TODO add created_at/updated_at
type Route struct {
	Id        int64      `db:"id"`
	Path      string     `db:"path"`
	Dest      string     `db:"dest"`
	Counter   int64      `db:"counter"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func initDB(file string) *gorp.DbMap {
	conn, err := sql.Open("sqlite3", file)
	checkErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: conn, Dialect: gorp.SqliteDialect{}}

	sql := `CREATE TABLE IF NOT EXISTS routes (
            id INTEGER NOT NULL PRIMARY KEY,
            path VARCHAR(255) NOT NULL UNIQUE,
            dest VARCHAR(255) NOT NULL,
            comment VARCHAR(255),
            counter INTEGER DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
          );`

	_, err = dbmap.Exec(sql)
	checkErr(err, "Create tables failed")

	dbmap.AddTableWithName(Route{}, "routes")

	return dbmap
}

func loadRoutes() []Route {
	var routes []Route
	_, err := db.Select(&routes, "SELECT id, path, dest FROM routes")
	checkErr(err, "Select failed")

	fmt.Printf("Loaded %d paths from database\n", len(routes))
	return routes
}

func incRouteCounter(route *Route) {
	_, err := db.Exec("UPDATE routes SET counter = counter + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", route.Id)
	checkErr(err, "Updating counter failed")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
