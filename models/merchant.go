package models

type Merchant struct {
	ID            string `json:"$id"`
	OwnerId       string `json:"owner_id"`
	OwnerName     string `json:"owner_name"`
	MerchantId    string `json:"merchant_id"`
	MerchantName  string `json:"merchant_name"`
	MerchantEmail string `json:"merchant_email"`
	Status        string `json:"status"`
	CreatedAt     string `json:"$createdAt"`
	UpdatedAt     string `json:"$updatedAt"`
}

type Cashier struct {
	ID           string `json:"$id"`
	MerchantId   string `json:"merchant_id"`
	CashierId    string `json:"cashier_id"`
	CashierName  string `json:"cashier_name"`
	CashierEmail string `json:"cashier_email"`
	Status       string `json:"status"`
	CreatedAt    string `json:"$createdAt"`
	UpdatedAt    string `json:"$updatedAt"`
}

const (
	StatusActive   = "active"
	StatusDeactive = "deactive"
)

func IsValidStatus(status string) bool {
	switch status {
	case StatusActive, StatusDeactive:
		return true
	default:
		return false
	}
}
