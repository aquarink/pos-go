package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func ProductList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		prods, err := client.ProductByUserID(os.Getenv("PRODUCTS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load products", http.StatusSeeOther)
			return
		}

		_, err = client.StoreByUserID(os.Getenv("STORES"), models.GlobalSessionData.UserId)
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/app/store?error=harap lengkapi profile toko anda terlebih dahulu", http.StatusSeeOther)
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
		category, err := client.CategoryByUserId(os.Getenv("CATEGORIES"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load products", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "Add New Product",
			Data: map[string]interface{}{
				"categories": category,
			},
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

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		if name == "" || category == "" || price == "" {
			http.Redirect(w, r, "/app/product/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		categoriesData, err := client.CategoryById(os.Getenv("CATEGORIES"), category)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=failed to load products", http.StatusSeeOther)
			return
		}

		priceInt, err := strconv.Atoi(price)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=invalid price", http.StatusSeeOther)
			return
		}

		check, _ := client.ProductByName(os.Getenv("PRODUCTS"), name, user_id)
		if check != nil {
			http.Redirect(w, r, "/app/product/add?error=nama produk sudah ada", http.StatusSeeOther)
			return
		}

		countPaket, err := client.StoreByUserID(os.Getenv("STORES"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=paket produk invalid", http.StatusSeeOther)
			return
		}

		maxProducts, err := strconv.Atoi(countPaket.Package[2])
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=invalid package limit", http.StatusSeeOther)
			return
		}

		prods, err := client.ProductByUserID(os.Getenv("PRODUCTS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=data produk invalid", http.StatusSeeOther)
			return
		}

		if len(prods) >= maxProducts {
			http.Redirect(w, r, "/app/product/add?error=anda tidak dapat menambah produk, harap upgrade paket", http.StatusSeeOther)
			return
		}

		// INI UPLOTAN
		file, _, err := r.FormFile("photo")
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to upload photo", http.StatusSeeOther)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp("", "")
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to create temp file", http.StatusSeeOther)
			return
		}
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to save temp file", http.StatusSeeOther)
			return
		}

		fileURL, fileID, fileNAME, err := client.FileUpload(os.Getenv("PRODUCTS_BUCKET"), tempFile.Name())
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to upload photo to server", http.StatusSeeOther)
			return
		}

		slug := utils.CreateSlug(name)

		now := time.Now().Format(time.RFC3339)

		product := models.Products{
			Name:      name,
			Category:  []string{categoriesData.ID, categoriesData.Name},
			Price:     priceInt,
			UserID:    user_id,
			Photo:     []string{fileURL, fileID, fileNAME, os.Getenv("APPWRITE_PROJECT_ID")},
			Slug:      slug,
			CreatedAt: now,
			UpdatedAt: now,
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

		category, err := client.CategoryByUserId(os.Getenv("CATEGORIES"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=failed to load products", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "Update Product",
			Data: map[string]interface{}{
				"product":    product,
				"categories": category,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/product/product_edit.html", data)
		return
	}
}

func ProductUpdate(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		category := r.FormValue("category")
		price := r.FormValue("price")
		productId := r.FormValue("productId")
		user_id := models.GlobalSessionData.UserId

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		if name == "" || category == "" || price == "" || productId == "" {
			http.Redirect(w, r, "/app/product/list?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		product, err := client.GetProductByID(os.Getenv("PRODUCTS"), productId)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=product not found", http.StatusSeeOther)
			return
		}

		categoriesData, err := client.CategoryById(os.Getenv("CATEGORIES"), category)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=failed to load categories", http.StatusSeeOther)
			return
		}

		priceInt, err := strconv.Atoi(price)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/app/product/edit/%s?error=invalid price", productId), http.StatusSeeOther)
			return
		}

		var fileURL string
		var fileID string
		var fileNAME string
		var projectID string

		file, _, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()

			if product != nil && len(product.Photo) > 0 {
				_ = client.FileRemove(os.Getenv("PRODUCTS_BUCKET"), product.Photo[1])
			}

			tempFile, err := os.CreateTemp("", "")
			if err != nil {
				http.Redirect(w, r, "/app/product/add?error=failed to create temp file", http.StatusSeeOther)
				return
			}
			defer tempFile.Close()
			defer os.Remove(tempFile.Name())

			_, err = io.Copy(tempFile, file)
			if err != nil {
				http.Redirect(w, r, "/app/product/add?error=failed to save temp file", http.StatusSeeOther)
				return
			}

			fileURL, fileID, fileNAME, err = client.FileUpload(os.Getenv("PRODUCTS_BUCKET"), tempFile.Name())
			if err != nil {
				http.Redirect(w, r, fmt.Sprintf("/app/product/edit/%s?error=failed to upload file", productId), http.StatusSeeOther)
				return
			}

			projectID = os.Getenv("APPWRITE_PROJECT_ID")
		} else {
			fileURL = product.Photo[0]
			fileID = product.Photo[1]
			fileNAME = product.Photo[2]
			projectID = product.Photo[1]
		}

		slug := utils.CreateSlug(name)

		now := time.Now().Format(time.RFC3339)

		productUpdate := models.Products{
			Name:      name,
			Category:  []string{categoriesData.ID, categoriesData.Name},
			Price:     priceInt,
			UserID:    user_id,
			Photo:     []string{fileURL, fileID, fileNAME, projectID},
			Slug:      slug,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = client.UpdateProduct(os.Getenv("PRODUCTS"), productId, productUpdate)
		if err != nil {
			http.Redirect(w, r, "/app/product/list?error=failed to create product", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/product/list?msg=product created successfully", http.StatusSeeOther)
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
