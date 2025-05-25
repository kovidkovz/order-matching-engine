package db

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
    var err error

    // 👇 Using 'kovid' as username
    DB, err = sql.Open("mysql", "kovid:Berzerk@27@tcp(localhost:3306)/ordersystem")
    if err != nil {
        log.Fatal("Failed to open DB:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Failed to ping DB:", err)
    }

    log.Println("✅ Connected to MariaDB!")
}
