package utils

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/rand"
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
		t, err = time.Parse("Jan 2, 2006, 15:04", date)
		if err != nil {
			return date
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

func CreateSlug(input string) string {
	slug := strings.ToLower(input)
	re := regexp.MustCompile(`[^\w\s]`)
	slug = re.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Trim(slug, "-")

	return slug
}

func GenerateUnique(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}
	// Generate random bytes
	randomBytes := make([]byte, length/2)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to hex
	randomID := hex.EncodeToString(randomBytes)

	// Ensure the ID is of the desired length
	if len(randomID) > length {
		randomID = randomID[:length]
	}

	return randomID, nil
}
