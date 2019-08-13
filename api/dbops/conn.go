package dbops

import (
	"database/sql"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	// dbConn, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/go_server?charset=utf8")
	dbConn, err = sql.Open("postgres", "postgres://db_user:123456@localhost/golang_db?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
}
