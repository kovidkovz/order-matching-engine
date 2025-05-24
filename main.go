package main

import (
	"github.com/gin-gonic/gin"
	"golang-order-matching-system/db"
	"golang-order-matching-system/api"
)

func main() {
	db.InitDB()

	router := gin.Default()
	router.POST("/orders", api.PlaceOrder)
	router.Run(":8080")
}
