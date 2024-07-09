package models

type Packages struct {
	ID               string `json:"$id"`
	Name             string `json:"name"`
	Price            int    `json:"price"`
	CashierAvailable int    `json:"cashier"`
	ProductAvailable int    `json:"product"`
	Description      string `json:"Description"`
	CreatedAt        string `json:"$createdAt"`
	UpdatedAt        string `json:"$updatedAt"`
}
