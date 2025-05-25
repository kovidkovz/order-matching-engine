package engine

import (
    "sort"
)

// Order represents a simplified order
type Order struct {
    ID       int
    Side     string  // "buy" or "sell"
    Price    float64
    Quantity int
    // You can add Timestamp or CreatedAt for time priority if needed
}

// Trade represents a matched trade
type Trade struct {
    BuyOrderID  int
    SellOrderID int
    Price       float64
    Quantity    int
}

// MatchOrders matches buy and sell orders and returns executed trades
func MatchOrders(buys []Order, sells []Order) ([]Trade, []Order, []Order) {
    trades := []Trade{}

    // Sort buy orders descending by price (highest first)
    sort.Slice(buys, func(i, j int) bool {
        return buys[i].Price > buys[j].Price
    })

    // Sort sell orders ascending by price (lowest first)
    sort.Slice(sells, func(i, j int) bool {
        return sells[i].Price < sells[j].Price
    })

    i, j := 0, 0
    for i < len(buys) && j < len(sells) {
        buy := &buys[i]
        sell := &sells[j]

        // If buy price < sell price, no match possible
        if buy.Price < sell.Price {
            break
        }

        // Calculate matched quantity
        qty := min(buy.Quantity, sell.Quantity)
        tradePrice := sell.Price // or some logic (like midpoint)

        trades = append(trades, Trade{
            BuyOrderID:  buy.ID,
            SellOrderID: sell.ID,
            Price:       tradePrice,
            Quantity:    qty,
        })

        // Deduct matched qty
        buy.Quantity -= qty
        sell.Quantity -= qty

        // Remove fully filled orders
        if buy.Quantity == 0 {
            i++
        }
        if sell.Quantity == 0 {
            j++
        }
    }

    // Remaining unmatched orders
    remainingBuys := buys[i:]
    remainingSells := sells[j:]

    return trades, remainingBuys, remainingSells
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
