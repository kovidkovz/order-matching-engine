package database

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
    var err error
    db, err = sql.Open("mysql", "kovid:Berzerk@27@tcp(127.0.0.1:3306)/ordersystem")
    if err != nil {
        log.Fatal("Error connecting to DB:", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal("Ping failed:", err)
    }
}

func GetDB() *sql.DB {
    return db
}
