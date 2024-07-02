package utils

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, layout string, tmpl string, data interface{}) {
	// parsedTemplate, err := template.ParseFiles(layout, tmpl)
	parsedTemplate, err := AddTemplateFuncs(template.New("layout")).ParseFiles(layout, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		dataMap[csrf.TemplateTag] = csrf.TemplateField(r)
	}

	err = parsedTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderTemplateWithSidebar(w http.ResponseWriter, r *http.Request, layout string, tmpl string, data interface{}) {
	// parsedTemplate, err := template.ParseFiles(layout, tmpl, "views/templates/sidebar.html")
	parsedTemplate, err := AddTemplateFuncs(template.New("layout")).ParseFiles(layout, tmpl, "views/templates/sidebar.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		dataMap[csrf.TemplateTag] = csrf.TemplateField(r)
	}

	err = parsedTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
