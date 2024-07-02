package utils

import (
	"html/template"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(str string) string {
	caser := cases.Title(language.English)
	return caser.String(str)
}

func DateFormat(date string, layout string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Format(layout)
}

func UcWords(str string) string {
	words := strings.Fields(str)
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, " ")
}

func AddTemplateFuncs(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{
		"title":      Title,
		"dateFormat": DateFormat,
		"ucwords":    UcWords,
	})
}
