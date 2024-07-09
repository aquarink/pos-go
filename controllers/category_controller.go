package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func CategoryList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		cat, err := client.ListCategory(os.Getenv("CATEGORIES"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=email atau password salah", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Categories",
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

func CategoryAdd(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
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

	if r.Method == http.MethodPost {
		name := r.FormValue("categoryName")
		user_id := models.GlobalSessionData.UserId

		if name == "" || user_id == "" {
			http.Redirect(w, r, "/app/category/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		check, _ := client.CategoryByNameAndUserId(os.Getenv("CATEGORIES"), name, user_id)
		if check != nil {
			http.Redirect(w, r, "/app/category/list?error=nama kategori sudah ada", http.StatusSeeOther)
			return
		}

		slugs := utils.CreateSlug(name)
		categoryData := models.Categories{
			Name:   name,
			Slug:   slugs,
			UserID: user_id,
		}

		err := client.CreateCategory(os.Getenv("CATEGORIES"), categoryData)
		if err != nil {
			http.Redirect(w, r, "/app/category/list?error=kesalahan data, harap coba kembali", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/category/list?msg=berhasil menambahkan kategori baru", http.StatusSeeOther)
	}
}

func CategoryEdit(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

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
}

func CategoryUpdate(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodPost {
		id := r.FormValue("categoryId")
		name := r.FormValue("categoryName")
		user_id := models.GlobalSessionData.UserId

		if id == "" || name == "" || user_id == "" {
			http.Redirect(w, r, "/app/category/list?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		slugs := utils.CreateSlug(name)
		category := models.Categories{
			ID:     id,
			Name:   name,
			Slug:   slugs,
			UserID: user_id,
		}

		_, err := client.UpdateCategory(os.Getenv("CATEGORIES"), id, category)
		if err != nil {
			http.Redirect(w, r, "/app/category/list?error=gagal edit kategori", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/category/list?msg=kategori berhasil di update", http.StatusSeeOther)
	}
}

func CategoryDelete(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/category/list?error=invalid data", http.StatusSeeOther)
			return
		}

		err := client.DeleteCategory(os.Getenv("CATEGORIES"), id)
		if err != nil {
			http.Redirect(w, r, "/app/category/list?error=kategori tidak ditemukan", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/category/list?msg=berhasil menghapus kategori", http.StatusSeeOther)
	}
}
