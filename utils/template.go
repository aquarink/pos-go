package utils

import (
	"html/template"
	"log"
	"net/http"
	"pos/models"

	"github.com/gorilla/csrf"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, layout string, tmpl string, data interface{}) {
	parsedTemplate, err := AddTemplateFuncs(template.New("layout")).ParseFiles(layout, tmpl)
	if err != nil {
		http.Redirect(w, r, "/app/signin?error=perbaikan grafik, harap coba lagi", http.StatusSeeOther)
		return
	}

	var dataMap map[string]interface{}
	if d, ok := data.(models.PublicData); ok {
		dataMap = map[string]interface{}{
			"Title":     d.Title,
			"Data":      d.Data,
			"Error":     d.Error,
			"Msg":       d.Msg,
			"Session":   d.Session,
			"CSRFToken": csrf.TemplateField(r),
		}
	} else {
		http.Redirect(w, r, "/app/signin?error=perbaikan sesi, harap coba lagi", http.StatusSeeOther)
		return
	}

	err = parsedTemplate.ExecuteTemplate(w, "layout", dataMap)
	if err != nil {
		http.Redirect(w, r, "/app/signin?error=perbaikan sistem tampilan, harap coba lagi", http.StatusSeeOther)
		return
	}
}

func RenderTemplateWithSidebar(w http.ResponseWriter, r *http.Request, layout string, tmpl string, data interface{}) {
	parsedTemplate, err := AddTemplateFuncs(template.New("layout")).ParseFiles(layout, tmpl, "views/templates/sidebar.html")
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/app/signin?error=perbaikan grafik, harap coba lagi", http.StatusSeeOther)
		return
	}

	var dataMap map[string]interface{}
	if d, ok := data.(models.PublicData); ok {
		dataMap = map[string]interface{}{
			"Title":     d.Title,
			"Data":      d.Data,
			"Error":     d.Error,
			"Msg":       d.Msg,
			"Session":   d.Session,
			"CSRFToken": csrf.TemplateField(r),
		}
	} else {
		http.Redirect(w, r, "/app/signin?error=perbaikan sesi, harap coba lagi", http.StatusSeeOther)
		return
	}

	err = parsedTemplate.ExecuteTemplate(w, "layout", dataMap)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/app/signin?error=perbaikan sistem tampilan, harap coba lagi", http.StatusSeeOther)
		return
	}
}
