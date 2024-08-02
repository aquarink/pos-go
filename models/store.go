package models

type Store struct {
	ID        string   `json:"$id"`
	User      []string `json:"user"` // user_id, name, email
	Name      string   `json:"name"`
	Address   []string `json:"address"` // city, address
	Logo      []string `json:"logo"`    // url, project_id
	Slug      string   `json:"slug"`
	Package   []string `json:"package"`
	Table     int      `json:"table"`
	CreatedAt string   `json:"$createdAt"`
	UpdatedAt string   `json:"$updatedAt"`
}

type Table struct {
	ID        string   `json:"$id"`
	UserId    string   `json:"user_id"`
	TableNo   int      `json:"table_no"`
	Code      string   `json:"code"`
	CodeImage []string `json:"code_image"`
	CreatedAt string   `json:"$createdAt"`
	UpdatedAt string   `json:"$updatedAt"`
}
