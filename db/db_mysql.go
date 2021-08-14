package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nikitamirzani323/gofiberapi/config"
)

var db *sql.DB
var err error

func Init() {
	conf := config.GetConfigMysql()
	connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic("koneksi string serror")
	}

	err = db.Ping()
	if err != nil {
		panic("DNS invalid")
	}
}
func CreateCon() *sql.DB {
	return db
}
