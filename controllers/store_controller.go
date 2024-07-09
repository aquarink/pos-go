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
	"strings"
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

		var logoURL string
		var projectID string

		file, _, err := r.FormFile("logo")
		if err == nil {
			defer file.Close()

			// RENAME
			uniqueFileNameBytes, err := exec.Command("uuidgen").Output()
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to generate unique file name", http.StatusSeeOther)
				return
			}

			uniqueFileName := strings.TrimSpace(string(uniqueFileNameBytes))

			tempFile, err := os.CreateTemp("", uniqueFileName)
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

			logoURL, err = client.UploadFile(os.Getenv("STORES_LOGO_BUCKET"), uniqueFileName, tempFile.Name())
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to upload file", http.StatusSeeOther)
				return
			}

			projectID = os.Getenv("APPWRITE_PROJECT_ID")
		} else {
			logoURL = stores.Logo[0]
			projectID = stores.Logo[1]
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
			Logo:      []string{logoURL, projectID},
			Slug:      slug,
			CreatedAt: created,
			UpdatedAt: now,
		}

		if stores != nil && stores.ID != "" {
			log.Println("UPDATE ------------")
			// UPDATE
			_, err = client.UpdateStore(os.Getenv("STORES"), user_id, updates)
			if err != nil {
				http.Redirect(w, r, "/app/store?error=failed to update", http.StatusSeeOther)
				return
			}
		} else {
			log.Println("CREATE ------------")
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
