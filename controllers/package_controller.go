package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func PackageList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {

		packages, err := client.ListPackage(os.Getenv("PACKAGES"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load package", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title:   "List of Packages",
			Data:    map[string]interface{}{"packages": packages},
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
		merchant := r.FormValue("merchant")
		cashier := r.FormValue("cashier")
		category := r.FormValue("category")
		product := r.FormValue("product")
		table := r.FormValue("table")
		description := r.FormValue("desc")

		check, _ := client.PackageByName(os.Getenv("PACKAGES"), name)
		if check != nil {
			http.Redirect(w, r, "/app/package/add?error=nama paket sudah ada", http.StatusSeeOther)
			return
		}

		priceInt, err := strconv.Atoi(price)
		if err != nil {
			http.Redirect(w, r, "/app/package/add?error=invalid price", http.StatusSeeOther)
			return
		}

		merchantInt, err := strconv.Atoi(merchant)
		if err != nil {
			http.Redirect(w, r, "/app/package/add?error=invalid merchant available", http.StatusSeeOther)
			return
		}

		cashierInt, err := strconv.Atoi(cashier)
		if err != nil {
			http.Redirect(w, r, "/app/package/add?error=invalid cashier available", http.StatusSeeOther)
			return
		}

		categoryInt, err := strconv.Atoi(category)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid category available", http.StatusSeeOther)
			return
		}

		productInt, err := strconv.Atoi(product)
		if err != nil {
			http.Redirect(w, r, "/app/package/add?error=invalid product available", http.StatusSeeOther)
			return
		}

		tableInt, err := strconv.Atoi(table)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid table available", http.StatusSeeOther)
			return
		}

		packageData := models.Packages{
			Name:              name,
			Price:             priceInt,
			MerchantAvailable: merchantInt,
			CashierAvailable:  cashierInt,
			CategoryAvailable: categoryInt,
			ProductAvailable:  productInt,
			TableAvailable:    tableInt,
			Description:       description,
		}

		err = client.CreatePackage(os.Getenv("PACKAGES"), packageData)
		if err != nil {
			http.Redirect(w, r, "/app/package/add?error=kesalahan data, harap coba kembali", http.StatusSeeOther)
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
		merchant := r.FormValue("merchant")
		cashier := r.FormValue("cashier")
		category := r.FormValue("category")
		product := r.FormValue("product")
		table := r.FormValue("table")
		description := r.FormValue("desc")

		if id == "" || name == "" || price == "" || merchant == "" || cashier == "" || category == "" || product == "" || table == "" || description == "" {
			http.Redirect(w, r, "/app/package/list?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		priceInt, err := strconv.Atoi(price)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid price", http.StatusSeeOther)
			return
		}

		cashierInt, err := strconv.Atoi(cashier)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid cashier available", http.StatusSeeOther)
			return
		}

		merchantInt, err := strconv.Atoi(merchant)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid merchant available", http.StatusSeeOther)
			return
		}

		categoryInt, err := strconv.Atoi(category)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid category available", http.StatusSeeOther)
			return
		}

		productInt, err := strconv.Atoi(product)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid product available", http.StatusSeeOther)
			return
		}

		tableInt, err := strconv.Atoi(table)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=invalid table available", http.StatusSeeOther)
			return
		}

		_, err = client.PackageById(os.Getenv("PACKAGES"), id)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=data tidak ditemukan", http.StatusSeeOther)
			return
		}

		packageData := models.Packages{
			Name:              name,
			Price:             priceInt,
			MerchantAvailable: merchantInt,
			CashierAvailable:  cashierInt,
			CategoryAvailable: categoryInt,
			ProductAvailable:  productInt,
			TableAvailable:    tableInt,
			Description:       description,
		}

		_, err = client.UpdatePackage(os.Getenv("PACKAGES"), id, packageData)
		if err != nil {
			http.Redirect(w, r, "/app/package/list?error=gagal edit paket", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/package/list?msg=paket berhasil di update", http.StatusSeeOther)
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
			http.Redirect(w, r, "/app/package/list?error=paket tidak ditemukan", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/package/list?msg=berhasil menghapus paket", http.StatusSeeOther)
	}
}
