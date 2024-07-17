package models

type User struct {
	ID            string `json:"$id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
	CreatedAt     string `json:"$createdAt"`
	UpdatedAt     string `json:"$updatedAt"`
}

const (
	RoleAdmin    = "admin"
	RoleMerchant = "merchant"
	RoleCashier  = "cashier"
	RoleKitchen  = "kitchen"
)

func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleMerchant, RoleCashier, RoleKitchen:
		return true
	default:
		return false
	}
}

type Cashier struct {
	ID           string `json:"$id"`
	MerchantId   string `json:"merchant_id"`
	CashierId    string `json:"cashier_id"`
	CashierName  string `json:"cashier_name"`
	CashierEmail string `json:"cashier_email"`
	CreatedAt    string `json:"$createdAt"`
	UpdatedAt    string `json:"$updatedAt"`
}

type Mails struct {
	ID        string `json:"$id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Text      string `json:"text"`
	HTML      string `json:"html"`
	CreatedAt string `json:"$createdAt"`
	UpdatedAt string `json:"$updatedAt"`
}
