package models

import "html/template"

type PublicData struct {
	Title     string
	Data      interface{}
	Error     string
	Msg       string
	Session   SessionData
	CSRFToken template.HTML
}

type SessionData struct {
	UserId string
	Role   string
}

// ngelempar data ke global session
var GlobalSessionData SessionData
