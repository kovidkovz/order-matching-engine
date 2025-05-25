package services

import (
    "database/sql"
    "log"
    "order_matching_engine/database"
    "order_matching_engine/models"
)

func MatchOrder(order models.Order) {
    conn := database.GetDB()
    tx, err := conn.Begin()
    if err != nil {
        log.Println("Error starting transaction:", err)
        return
    }

    query := `SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at
              FROM orders WHERE symbol = ? AND side = ? AND price <= ? AND status = 'open'
              ORDER BY price ASC, created_at ASC`

    var rows *sql.Rows
    if order.Side == "buy" {
        rows, err = tx.Query(query, order.Symbol, "sell", order.Price)
    } else {
        query = `SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at
                 FROM orders WHERE symbol = ? AND side = ? AND price >= ? AND status = 'open'
                 ORDER BY price DESC, created_at ASC`
        rows, err = tx.Query(query, order.Symbol, "buy", order.Price)
    }

    if err != nil {
        tx.Rollback()
        return
    }
    defer rows.Close()

    for rows.Next() {
        var match models.Order
        err := rows.Scan(&match.ID, &match.Symbol, &match.Side, &match.Type, &match.Price,
            &match.Quantity, &match.RemainingQuantity, &match.Status, &match.CreatedAt)
        if err != nil {
            tx.Rollback()
            return
        }

        matchQty := min(order.RemainingQuantity, match.RemainingQuantity)

        _, err = tx.Exec(`INSERT INTO trades (buy_order_id, sell_order_id, price, quantity)
                          VALUES (?, ?, ?, ?)`,
            getBuyer(order, match).ID, getSeller(order, match).ID, match.Price, matchQty)
        if err != nil {
            tx.Rollback()
            return
        }

        order.RemainingQuantity -= matchQty
        match.RemainingQuantity -= matchQty

        if match.RemainingQuantity == 0 {
            tx.Exec("UPDATE orders SET status='filled' WHERE id=?", match.ID)
        } else {
            tx.Exec("UPDATE orders SET remaining_quantity=? WHERE id=?", match.RemainingQuantity, match.ID)
        }

        if order.RemainingQuantity == 0 {
            tx.Exec("UPDATE orders SET status='filled' WHERE id=?", order.ID)
            break
        }
    }

    if order.RemainingQuantity > 0 {
        tx.Exec("UPDATE orders SET remaining_quantity=? WHERE id=?", order.RemainingQuantity, order.ID)
    }

    tx.Commit()
}

func getBuyer(a, b models.Order) models.Order {
    if a.Side == "buy" {
        return a
    }
    return b
}

func getSeller(a, b models.Order) models.Order {
    if a.Side == "sell" {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
