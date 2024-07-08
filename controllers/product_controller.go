package controllers

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"pos/models"
	"pos/services"
	"pos/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func ProductList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		prods, err := client.ListProducts(os.Getenv("PRODUCTS"))
		if err != nil {
			log.Println(err.Error())
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

		if name == "" || category == "" || price == "" || user_id == "" {
			http.Redirect(w, r, "/app/product/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		categoriesData, err := client.CategoryById(os.Getenv("CATEGORIES"), category)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load products", http.StatusSeeOther)
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

		// Generate a unique file name
		uniqueFileNameBytes, err := exec.Command("uuidgen").Output()
		if err != nil {
			http.Redirect(w, r, "/app/product/add?error=failed to generate unique file name", http.StatusSeeOther)
			return
		}

		uniqueFileName := strings.TrimSpace(string(uniqueFileNameBytes))

		tempFile, err := os.CreateTemp("", uniqueFileName)
		if err != nil {

			log.Println(err.Error())
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

		photoURL, err := client.UploadFile(os.Getenv("PRODUCTS_BUCKET"), uniqueFileName, tempFile.Name())
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/app/product/add?error=failed to upload photo to server", http.StatusSeeOther)
			return
		}

		slug := utils.CreateSlug(name)

		now := time.Now().Format(time.RFC3339)

		log.Println("uniqueFileName : " + uniqueFileName)
		log.Println("resPhoto : " + photoURL)
		product := models.Products{
			Name:      name,
			Category:  []string{categoriesData.ID, categoriesData.Name},
			Price:     priceInt,
			UserID:    user_id,
			Photo:     []string{photoURL, os.Getenv("APPWRITE_PROJECT_ID")},
			Slug:      slug,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = client.CreateProduct(os.Getenv("PRODUCTS"), product)
		if err != nil {
			log.Printf(err.Error())
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

		categoriesData, err := client.CategoryById(os.Getenv("CATEGORIES"), category)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load products", http.StatusSeeOther)
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
			Category: []string{categoriesData.ID, categoriesData.Name},
			Price:    priceInt,
			UserID:   user_id,
			Slug:     slug,
		}

		if photoID != "" {
			product.Photo = []string{photoID, os.Getenv("APPWRITE_PROJECT_ID")}
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
