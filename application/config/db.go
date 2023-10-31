package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB  *sql.DB
	err error
)

func init() {
	//dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
	//	"user", "password", "127.0.0.1:3306", "coolname_db")
	DB, err = sql.Open("mysql", "root:password@tcp(localhost:33062)/recycle_db")

	if err != nil {
		panic(err.Error())
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}

	log.Println(" Connected to Database ")

}
