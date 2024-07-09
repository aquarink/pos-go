package models

type Store struct {
	ID        string   `json:"$id"`
	User      []string `json:"user"` // user_id, name, email
	Name      string   `json:"name"`
	Address   []string `json:"address"` // city, address
	Logo      []string `json:"logo"`    // url, project_id
	Slug      string   `json:"slug"`
	CreatedAt string   `json:"$createdAt"`
	UpdatedAt string   `json:"$updatedAt"`
}
