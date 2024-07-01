package models

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type Mails struct {
	ID      string
	UserID  string
	Email   string
	Subject string
	Text    string
	HTML    string
}
