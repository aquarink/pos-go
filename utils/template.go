package utils

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, layout string, tmpl string, data interface{}) {
	// parsedTemplate, err := template.ParseFiles(layout, tmpl)
	parsedTemplate, err := AddTemplateFuncs(template.New("layout")).ParseFiles(layout, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = parsedTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderTemplateWithSidebar(w http.ResponseWriter, layout string, tmpl string, data interface{}) {
	// parsedTemplate, err := template.ParseFiles(layout, tmpl, "views/templates/sidebar.html")
	parsedTemplate, err := AddTemplateFuncs(template.New("layout")).ParseFiles(layout, tmpl, "views/templates/sidebar.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = parsedTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
