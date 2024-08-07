package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func AddTemplateFuncs(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{
		"title":        Title,
		"dateFormat":   DateFormat,
		"ucwords":      UcWords,
		"coma":         Comma,
		"splits":       Splits,
		"viewDownload": ReplaceViewWithDownload,
	})
}

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

func CreateSlug(input string) string {
	slug := strings.ToLower(input)
	re := regexp.MustCompile(`[^\w\s]`)
	slug = re.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Trim(slug, "-")

	return slug
}

func Comma(value interface{}) string {
	switch v := value.(type) {
	case int:
		if v == 0 {
			return "0"
		}
		return humanize.Comma(int64(v))
	case float64:
		if v == 0.0 {
			return "0"
		}
		return humanize.Commaf(v)
	case float32:
		if v == 0.0 {
			return "0"
		}
		return humanize.Commaf(float64(v))
	case string:
		// Attempt to parse string to float
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			if f == 0.0 {
				return "0"
			}
			return humanize.Commaf(f)
		}
		// Attempt to parse string to int
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			if i == 0 {
				return "0"
			}
			return humanize.Comma(i)
		}
		// If parsing fails, return original string or "0" if empty
		if v == "" {
			return "0"
		}
		return v
	default:
		// If the type is not handled, return "0"
		return "0"
	}
}

func Uniqid(moreEntropy bool) string {
	sec := time.Now().Unix()
	uniqid := fmt.Sprintf("%x", sec)

	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		panic(err)
	}
	uniqid += hex.EncodeToString(randBytes)

	if moreEntropy {
		extra, _ := rand.Int(rand.Reader, big.NewInt(100000))
		uniqid += fmt.Sprintf("%08d", extra)
	}

	return uniqid
}

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
}

func Splits(s, sep string) []string {
	return strings.Split(s, sep)
}

func ReplaceViewWithDownload(url string) string {
	return strings.Replace(url, "view", "download", 1)
}
