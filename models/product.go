package models

type Categories struct {
	ID        string `json:"$id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"$createdAt"`
	UpdatedAt string `json:"$updatedAt"`
}
