package utils

import (
	"html/template"
	"regexp"
	"strconv"
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

func NumberFormat(number float64, decimals int) string {
	decPoint := "."
	thousandsSep := ","

	// ANGKA KE STRING
	strNumber := strconv.FormatFloat(number, 'f', decimals, 64)

	// SPLIT
	parts := strings.Split(strNumber, ".")
	integerPart := parts[0]
	fractionalPart := ""
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	// KASIH PER NOL 3
	result := ""
	count := 0
	for i := len(integerPart) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result = thousandsSep + result
		}
		result = string(integerPart[i]) + result
		count++
	}

	if decimals > 0 {
		result = result + decPoint + fractionalPart
	}

	return result
}
