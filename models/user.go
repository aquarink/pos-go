package models

type User struct {
	ID            string `json:"$id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
}

// Define constants for valid roles
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

type Mails struct {
	ID      string `json:"$id"`
	UserID  string
	Email   string
	Subject string
	Text    string
	HTML    string
}
