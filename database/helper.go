package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "amartha"
	dbName   = "amartha-db"
)

// ConnectDB provides a connection to mysql db. It is the caller's responsibility to close the connection.
func ConnectDB() *sqlx.DB {
	var connectString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := sqlx.Connect("mysql", connectString)
	if err != nil {
		panic(err)
	}

	return db
}
