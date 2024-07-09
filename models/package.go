package models

type Packages struct {
	ID               string `json:"$id"`
	Name             string `json:"name"`
	Price            string `json:"price"`
	CashierAvailable string `json:"cashier"`
	ProductAvailable string `json:"product"`
	Description      string `json:"description"`
	CreatedAt        string `json:"$createdAt"`
	UpdatedAt        string `json:"$updatedAt"`
}
