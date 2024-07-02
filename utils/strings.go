package utils

import (
	"html/template"
	"log"
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
	log.Printf("Original date: %s", date) // Log the original date string

	// Try to parse the date using ISO 8601 format
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// If parsing fails, assume the date is already in the desired format
		log.Printf("Error parsing date with RFC3339: %v", err)

		// Attempt parsing with another common format
		t, err = time.Parse("Jan 2, 2006, 15:04", date)
		if err != nil {
			log.Printf("Error parsing date with custom format: %v", err)
			return date // Return the original string if parsing fails
		}
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
