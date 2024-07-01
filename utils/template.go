package utils

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, layout string, tmpl string, data interface{}) {
	log.Println("Parsing templates:", layout, tmpl)
	parsedTemplate, err := template.ParseFiles(layout, tmpl)
	if err != nil {
		log.Println("Error parsing templates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Executing template with layout")
	err = parsedTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderTemplateWithSidebar(w http.ResponseWriter, layout string, tmpl string, data interface{}) {
	log.Println("Parsing templates:", layout, tmpl)
	parsedTemplate, err := template.ParseFiles(layout, tmpl, "views/templates/sidebar.html")
	if err != nil {
		log.Println("Error parsing templates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Executing template with layout")
	err = parsedTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
