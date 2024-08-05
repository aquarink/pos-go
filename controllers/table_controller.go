package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/sessions"
)

func TableList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		tab, err := client.ListTables(os.Getenv("TABLES"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=table data error", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "QR Tables",
			Data: map[string]interface{}{
				"tables": tab,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/table.html", data)
		return
	}
}

func TableNoGenerate(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {

		name := r.FormValue("categoryName")
		user_id := models.GlobalSessionData.UserId

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		if name == "" {
			http.Redirect(w, r, "/app/chair?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		check, _ := client.CategoryByNameAndUserId(os.Getenv("CATEGORIES"), name, user_id)
		if check != nil {
			http.Redirect(w, r, "/app/chair?error=nama kategori sudah ada", http.StatusSeeOther)
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
			http.Redirect(w, r, "/app/chair?error=kesalahan data, harap coba kembali", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/chair?msg=berhasil menambahkan kategori baru", http.StatusSeeOther)
	}
}
