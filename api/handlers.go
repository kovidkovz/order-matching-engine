package api

import (
    "encoding/json"
    "net/http"
    "order_matching_engine/database"
    "order_matching_engine/models"
    "order_matching_engine/services"
)

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
    var order models.Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    conn := database.GetDB()
    res, err := conn.Exec(`
        INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
        order.Symbol, order.Side, order.Type, order.Price, order.Quantity, order.Quantity, "open")
    if err != nil {
        http.Error(w, "Failed to insert order", http.StatusInternalServerError)
        return
    }

    id, _ := res.LastInsertId()
    order.ID = int(id)
    order.RemainingQuantity = order.Quantity
    order.Status = "open"

    services.MatchOrder(order)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}
