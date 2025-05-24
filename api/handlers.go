package api

import (
    "net/http"
    "order_matching_engine/db"

    "github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
    // parse input, insert order using raw SQL in db.DB
    c.JSON(http.StatusOK, gin.H{"message": "order created"})
}

func ListOrders(c *gin.Context) {
    // query orders from db.DB and return as JSON
    c.JSON(http.StatusOK, gin.H{"orders": []string{}})
}
