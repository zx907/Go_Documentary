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
	dbConn, err = sql.Open("postgres", "db_user:123456@tcp(localhost:5432)/golang_db?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
}
