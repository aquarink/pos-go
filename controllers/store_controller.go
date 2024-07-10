package controllers

import (
	"io"
	"log"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

func StoreEdit(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		user_id := models.GlobalSessionData.UserId

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		stores, _ := client.StoreByUserID(os.Getenv("STORES"), user_id)

		if stores == nil {
			stores = &models.Store{}
		}

		data := models.PublicData{
			Title: "Store Profile",
			Data: map[string]interface{}{
				"stores": stores,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/store.html", data)
		return
	}
}

func StoreUpdate(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		city := r.FormValue("city")
		address := r.FormValue("address")
		user_id := models.GlobalSessionData.UserId

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		if name == "" || city == "" || address == "" {
			http.Redirect(w, r, "/app/store?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		stores, _ := client.StoreByUserID(os.Getenv("STORES"), user_id)

		user, err := client.GetUserByID(os.Getenv("USERS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/dahboard?error=data anda tidak valid", http.StatusSeeOther)
			return
		}

		var fileURL string
		var fileID string
		var fileNAME string
		var projectID string

		file, _, err := r.FormFile("logo")
		if err == nil {
			defer file.Close()

			if stores != nil && len(stores.Logo) > 0 {
				_ = client.FileRemove(os.Getenv("STORES_LOGO_BUCKET"), stores.Logo[1])
			}

			tempFile, err := os.CreateTemp("", "")
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to create temp file", http.StatusSeeOther)
				return
			}
			defer tempFile.Close()
			defer os.Remove(tempFile.Name())

			_, err = io.Copy(tempFile, file)
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to save temp file", http.StatusSeeOther)
				return
			}

			fileURL, fileID, fileNAME, err = client.FileUpload(os.Getenv("STORES_LOGO_BUCKET"), tempFile.Name())
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to upload file", http.StatusSeeOther)
				return
			}

			projectID = os.Getenv("APPWRITE_PROJECT_ID")
		} else {
			fileURL = stores.Logo[0]
			fileID = stores.Logo[1]
			fileNAME = stores.Logo[2]
			projectID = stores.Logo[3]
		}

		slug := utils.CreateSlug(name)
		now := time.Now().Format(time.RFC3339)
		created := now

		// PACKAGE
		packageName := "free"
		packageCashier := 1
		packageProduct := 2

		if stores != nil {
			created = stores.CreatedAt

			// package
			if len(stores.Package) >= 3 {
				if stores.Package[0] != "" {
					packageName = stores.Package[0]
				}
				if stores.Package[1] != "" {
					packageCashier, err = strconv.Atoi(stores.Package[1])
					if err != nil {
						log.Println(err.Error())
					}
				}
				if stores.Package[2] != "" {
					packageProduct, err = strconv.Atoi(stores.Package[2])
					if err != nil {
						log.Println(err.Error())
					}
				}
			}
		}

		updates := models.Store{
			User:      []string{user.ID, user.Name},
			Name:      name,
			Address:   []string{city, address},
			Logo:      []string{fileURL, fileID, fileNAME, projectID},
			Slug:      slug,
			Package:   []string{packageName, strconv.Itoa(packageCashier), strconv.Itoa(packageProduct)},
			CreatedAt: created,
			UpdatedAt: now,
		}

		if stores != nil && stores.ID != "" {
			// UPDATE
			_, err = client.UpdateStore(os.Getenv("STORES"), user_id, updates)
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to update", http.StatusSeeOther)
				return
			}
		} else {
			// ADD
			err = client.CreateStore(os.Getenv("STORES"), updates)
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to create", http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/app/store?msg=product created successfully", http.StatusSeeOther)
	}
}
