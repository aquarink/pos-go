package controllers

import (
	"io"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func ProductList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		prods, err := client.ListProducts(os.Getenv("PRODUCTS"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load products", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Products",
			Data: map[string]interface{}{
				"products": prods,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/product_list.html", data)
		return
	}
}

func ProductAdd(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		data := models.PublicData{
			Title:   "Add New Product",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/product_add.html", data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		category := r.FormValue("category")
		price := r.FormValue("price")
		user_id := models.GlobalSessionData.UserId

		if name == "" || category == "" || price == "" || user_id == "" {
			http.Redirect(w, r, "/app/product/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		priceInt, err := strconv.Atoi(price)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=invalid price", http.StatusSeeOther)
			return
		}

		// Handle file upload
		file, _, err := r.FormFile("photo")
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to upload photo", http.StatusSeeOther)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp("", "upload-*.png")
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to create temp file", http.StatusSeeOther)
			return
		}
		defer os.Remove(tempFile.Name())

		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to save temp file", http.StatusSeeOther)
			return
		}

		photoID, err := client.UploadFile(os.Getenv("BUCKET_ID"), "unique()", tempFile.Name())
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to upload photo to server", http.StatusSeeOther)
			return
		}

		slug := utils.CreateSlug(name)

		product := models.Products{
			Name:     name,
			Category: category,
			Price:    priceInt,
			UserID:   user_id,
			Photo:    photoID,
			Slug:     slug,
		}

		err = client.CreateProduct(os.Getenv("PRODUCTS"), product)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to create product", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/product/list?msg=product created successfully", http.StatusSeeOther)
	}
}

func ProductEdit(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/product/list?error=invalid data", http.StatusSeeOther)
			return
		}

		product, err := client.GetProductByID(os.Getenv("PRODUCTS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=product not found", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title:   "Edit Product",
			Data:    map[string]interface{}{"product": product},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/product_edit.html", data)
		return
	}

	if r.Method == http.MethodPost {
		id := r.FormValue("productId")
		name := r.FormValue("name")
		category := r.FormValue("category")
		price := r.FormValue("price")
		user_id := models.GlobalSessionData.UserId

		if id == "" || name == "" || category == "" || price == "" || user_id == "" {
			http.Redirect(w, r, "/app/product/edit/"+id+"?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		priceInt, err := strconv.Atoi(price)
		if err != nil {
			http.Redirect(w, r, "/app/product/edit/"+id+"?error=invalid price", http.StatusSeeOther)
			return
		}

		// Handle file upload
		var photoID string
		file, _, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()

			tempFile, err := os.CreateTemp("", "upload-*.png")
			if err != nil {
				http.Redirect(w, r, "/app/product/edit/"+id+"?error=failed to create temp file", http.StatusSeeOther)
				return
			}
			defer os.Remove(tempFile.Name())

			_, err = io.Copy(tempFile, file)
			if err != nil {
				http.Redirect(w, r, "/app/product/edit/"+id+"?error=failed to save temp file", http.StatusSeeOther)
				return
			}

			photoID, err = client.UploadFile(os.Getenv("BUCKET_ID"), "unique()", tempFile.Name())
			if err != nil {
				http.Redirect(w, r, "/app/product/edit/"+id+"?error=failed to upload photo to server", http.StatusSeeOther)
				return
			}
		}

		slug := utils.CreateSlug(name)

		product := models.Products{
			ID:       id,
			Name:     name,
			Category: category,
			Price:    priceInt,
			UserID:   user_id,
			Slug:     slug,
		}

		if photoID != "" {
			product.Photo = photoID
		}

		_, err = client.UpdateProduct(os.Getenv("PRODUCTS"), id, product)
		if err != nil {
			http.Redirect(w, r, "/app/product/edit/"+id+"?error=failed to update product", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/product/list?msg=product updated successfully", http.StatusSeeOther)
	}
}

func ProductDelete(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/product/list?error=invalid data", http.StatusSeeOther)
			return
		}

		err := client.DeleteProduct(os.Getenv("PRODUCTS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=product not found", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/product/list?msg=product deleted successfully", http.StatusSeeOther)
	}
}
