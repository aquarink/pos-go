package models

type PublicData struct {
	Title   string
	Data    interface{}
	Error   string
	Msg     string
	Session SessionData
}

type SessionData struct {
	Role      string
	LastLogin string
}

var GlobalSessionData SessionData
