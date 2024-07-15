package models

type Checkout struct {
	ID            string   `json:"$id"`
	UserId        string   `json:"user_id"`
	Queue         string   `json:"queue"`
	TrxId         string   `json:"trx_id"`
	DineType      string   `json:"dine_type"`
	TableNumber   string   `json:"table_number"`
	Items         []string `json:"items"` // productIDs, productNames, productPrices, productQtys, productNotes
	TotalItem     int      `json:"total_item"`
	Tax           float64  `json:"tax"`
	TaxTotal      float64  `json:"tax_total"`
	TotalPayment  float64  `json:"total_payment"`
	PaymentMethod string   `json:"payment_method"`
	Change        float64  `json:"change"`
	Createdat     string   `json:"created_at"`
}
