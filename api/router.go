package api

import (
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.POST("/order", CreateOrder)
    r.GET("/orders", ListOrders)
    return r
}
