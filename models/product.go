package models

type Categories struct {
	ID        string `json:"$id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"$createdAt"`
	UpdatedAt string `json:"$updatedAt"`
}

type Products struct {
	ID        string   `json:"$id"`
	Name      string   `json:"name"`
	Category  []string `json:"category"`
	Price     int      `json:"price"`
	UserID    string   `json:"user_id"`
	Photo     []string `json:"photo"`
	Slug      string   `json:"slug"`
	CreatedAt string   `json:"$createdAt"`
	UpdatedAt string   `json:"$updatedAt"`
}
