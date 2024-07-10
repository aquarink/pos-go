package controllers

import (
	"io"
	"log"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
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

		stores, _ := client.GetStoreByUserID(os.Getenv("STORES"), user_id)

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

		stores, _ := client.GetStoreByUserID(os.Getenv("STORES"), user_id)

		user, err := client.GetUserByID(os.Getenv("USERS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/dahboard?error=data anda tidak valid", http.StatusSeeOther)
			return
		}

		var fileURL string
		var fileID string
		var fileNAME string
		var projectID string

		// uniqueFileNameBytes, err := exec.Command("uuidgen").Output()
		// if err != nil {
		// 	http.Redirect(w, r, "/app/store?error=failed to generate unique file name", http.StatusSeeOther)
		// 	return
		// }

		// uniqueFileName := strings.TrimSpace(string(uniqueFileNameBytes))

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
				log.Println(err.Error())
				http.Redirect(w, r, "/app/store?error=failed to upload file", http.StatusSeeOther)
				return
			}

			log.Println(fileURL)
			log.Println(fileID)
			log.Println(fileNAME)

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

		if stores != nil {
			created = stores.CreatedAt
		}

		updates := models.Store{
			User:      []string{user.ID, user.Name},
			Name:      name,
			Address:   []string{city, address},
			Logo:      []string{fileURL, fileID, fileNAME, projectID},
			Slug:      slug,
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
