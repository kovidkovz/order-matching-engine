package models

type Order struct {
	ID               int     `json:"id"`
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"`   // buy or sell
	Type             string  `json:"type"`   // limit or market
	Price            float64 `json:"price"`
	Quantity         int     `json:"quantity"`
	RemainingQuantity int    `json:"remaining_quantity"`
	Status           string  `json:"status"` // open, filled, canceled
}
