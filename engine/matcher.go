package engine

import (
	"database/sql"
	"log"
	"golang-order-matching-system/database"
	"golang-order-matching-system/models"
)

func MatchOrder(order models.Order) {
	tx, err := db.DB.Begin()
	if err != nil {
		log.Println("Transaction begin failed:", err)
		return
	}
	defer tx.Commit()

	var rows *sql.Rows
	if order.Side == "buy" {
		rows, err = tx.Query(`SELECT id, price, remaining_quantity FROM orders 
			WHERE symbol=? AND side='sell' AND status='open' 
			AND (type='limit') ORDER BY price ASC, created_at ASC`, order.Symbol)
	} else {
		rows, err = tx.Query(`SELECT id, price, remaining_quantity FROM orders 
			WHERE symbol=? AND side='buy' AND status='open' 
			AND (type='limit') ORDER BY price DESC, created_at ASC`, order.Symbol)
	}
	if err != nil {
		log.Println("Query failed:", err)
		return
	}
	defer rows.Close()

	remaining := order.Quantity

	for rows.Next() && remaining > 0 {
		var matchID int
		var matchPrice float64
		var matchQty int
		rows.Scan(&matchID, &matchPrice, &matchQty)

		if order.Type == "limit" && ((order.Side == "buy" && order.Price < matchPrice) || (order.Side == "sell" && order.Price > matchPrice)) {
			break
		}

		tradeQty := min(remaining, matchQty)

		// Insert into trades
		if order.Side == "buy" {
			tx.Exec("INSERT INTO trades (buy_order_id, sell_order_id, symbol, price, quantity) VALUES (?, ?, ?, ?, ?)",
				order.ID, matchID, order.Symbol, matchPrice, tradeQty)
		} else {
			tx.Exec("INSERT INTO trades (buy_order_id, sell_order_id, symbol, price, quantity) VALUES (?, ?, ?, ?, ?)",
				matchID, order.ID, order.Symbol, matchPrice, tradeQty)
		}

		// Update matched order
		tx.Exec("UPDATE orders SET remaining_quantity = remaining_quantity - ?, status = CASE WHEN remaining_quantity - ? = 0 THEN 'filled' ELSE 'open' END WHERE id = ?",
			tradeQty, tradeQty, matchID)

		remaining -= tradeQty
	}

	status := "open"
	if remaining == 0 {
		status = "filled"
	} else if order.Type == "market" {
		status = "canceled"
	}

	tx.Exec("INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		order.Symbol, order.Side, order.Type, order.Price, order.Quantity, remaining, status)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
