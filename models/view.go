package models

type PublicData struct {
	Title   string
	Data    interface{}
	Error   string
	Msg     string
	Session SessionData
}

type SessionData struct {
	UserId    string
	Role      string
	LastLogin string
}

// ngelempar data ke global session
var GlobalSessionData SessionData
