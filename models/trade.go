package models

import "time"

type Trade struct {
    ID         int       `json:"id"`
    BuyOrderID int       `json:"buy_order_id"`
    SellOrderID int      `json:"sell_order_id"`
    Price      float64   `json:"price"`
    Quantity   int       `json:"quantity"`
    CreatedAt  time.Time `json:"created_at"`
}
