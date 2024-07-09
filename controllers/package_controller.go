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

func PackageList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		data := models.PublicData{
			Title:   "List of Category",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/package/package_list.html", data)
		return
	}
}

func PackageAdd(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {

		data := models.PublicData{
			Title:   "Add New Package",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/package/package_add.html", data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		price := r.FormValue("price")
		cashier := r.FormValue("cashier")
		product := r.FormValue("product")
		description := r.FormValue("desc")

		packageData := models.Packages{
			Name:             name,
			Price:            price,
			CashierAvailable: cashier,
			ProductAvailable: product,
			Description:      description,
		}

		err := client.CreatePackage(os.Getenv("PACKAGES"), packageData)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=kesalahan data, harap coba kembali", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/package/list?msg=berhasil menambahkan paket baru", http.StatusSeeOther)
	}
}

func PackageEdit(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/package/list?error=invalid data", http.StatusSeeOther)
			return
		}

		packages, err := client.PackageById(os.Getenv("PACKAGES"), id)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=category not found", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title:   "Edit Package",
			Data:    map[string]interface{}{"package": packages},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/package/package_edit.html", data)
		return
	}
}

func PackageUpdate(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		price := r.FormValue("price")
		cashier := r.FormValue("cashier")
		product := r.FormValue("product")
		description := r.FormValue("desc")

		if id == "" || name == "" || price == "" || cashier == "" || product == "" || description == "" {
			http.Redirect(w, r, "/app/package/list?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		packageData := models.Packages{
			Name:             name,
			Price:            price,
			CashierAvailable: cashier,
			ProductAvailable: product,
			Description:      description,
		}

		_, err := client.UpdatePackage(os.Getenv("PACKAGES"), id, packageData)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=gagal edit kategori", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/package/list?msg=kategori berhasil di update", http.StatusSeeOther)
	}
}

func PackageDelete(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/package/list?error=invalid data", http.StatusSeeOther)
			return
		}

		err := client.DeletePackage(os.Getenv("PACKAGES"), id)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=kategori tidak ditemukan", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/package/list?msg=berhasil menghapus kategori", http.StatusSeeOther)
	}
}
