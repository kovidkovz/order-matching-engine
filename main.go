package main

import (
    "log"
    "net/http"
    "order_matching_engine/api"
    "order_matching_engine/database"
)

func main() {
    database.InitDB()

    http.HandleFunc("/orders", api.PlaceOrder)

    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
