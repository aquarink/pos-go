package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
)

func CategoryList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient) {
	if r.Method == http.MethodGet {
		cat, err := client.ListCategory(os.Getenv("CATEGORIES"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=email atau password salah", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Category",
			Data: map[string]interface{}{
				"categories": cat,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/category_list.html", data)
		return
	}
}

func CategoryAdd(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient) {
	if r.Method == http.MethodGet {

		data := models.PublicData{
			Title:   "Add New Category",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/category_add.html", data)
		return
	}

	//
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		user_id := models.GlobalSessionData.UserId

		if name == "" || user_id == "" {
			http.Redirect(w, r, "/app/category/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}
	}
}

func CategoryEdit(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient) {
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("data")

		if id == "" {
			http.Redirect(w, r, "/app/category/list?error=invalid data", http.StatusSeeOther)
			return
		}

		category, err := client.CategoryById(os.Getenv("CATEGORIES"), id)
		if err != nil {
			http.Redirect(w, r, "/app/category/list?error=category not found", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title:   "Edit Category",
			Data:    map[string]interface{}{"category": category},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/category_edit.html", data)
		return
	}

	//
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		user_id := models.GlobalSessionData.UserId

		if id == "" || name == "" || user_id == "" {
			http.Redirect(w, r, "/app/category/edit?data="+id+"&error=form tidak lengkap", http.StatusSeeOther)
			return
		}
	}
}
