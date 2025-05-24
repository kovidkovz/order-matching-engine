package main

import (
    "order_matching_engine/database"
    "order_matching_engine/api"
)

func main() {
    db.InitDB()
    router := api.SetupRouter()
    router.Run(":8080")
}
