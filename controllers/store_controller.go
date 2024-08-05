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

		// OWNER
		// merchant di merchants by cashier. ambil MerchantName
		merchantData, err := client.MerchantByMerchantId(os.Getenv("MERCHANTS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed merchant data", http.StatusSeeOther)
			return
		}

		//
		ownerData, err := client.OwnerDataByOwnerId(os.Getenv("OWNERS"), merchantData[0].OwnerId)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed owner owner", http.StatusSeeOther)
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
				"owner":  ownerData,
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
		tabl := r.FormValue("table")
		user_id := models.GlobalSessionData.UserId

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		if name == "" || city == "" || address == "" || tabl == "" {
			http.Redirect(w, r, "/app/store?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		stores, _ := client.StoreByUserID(os.Getenv("STORES"), user_id)

		// merchant di merchants by cashier. ambil MerchantName
		merchantData, err := client.MerchantByMerchantId(os.Getenv("MERCHANTS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/store?error=failed merchant data", http.StatusSeeOther)
			return
		}

		//
		ownerData, err := client.OwnerDataByOwnerId(os.Getenv("OWNERS"), merchantData[0].OwnerId)
		if err != nil {
			http.Redirect(w, r, "/app/store?error=failed owner owner", http.StatusSeeOther)
			return
		}

		// INI UPLOTAN
		var fileURL, fileID, fileNAME, projectID string

		// Cek apakah ada file logo yang diupload
		file, _, err := r.FormFile("logo")
		if err == nil {
			defer file.Close()

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

			fileURL, fileID, fileNAME, projectID, err = client.FileUpload(os.Getenv("STORES_LOGO_BUCKET"), tempFile.Name())
			if err != nil {
				if stores != nil && len(stores.Logo) > 0 {
					fileURL = stores.Logo[0]
					fileID = stores.Logo[1]
					fileNAME = stores.Logo[2]
					projectID = stores.Logo[3]
				} else {
					http.Redirect(w, r, "/app/store?error=failed to upload file", http.StatusSeeOther)
					return
				}
			} else {
				if stores != nil && len(stores.Logo) > 0 {
					_ = client.FileRemove(os.Getenv("STORES_LOGO_BUCKET"), stores.Logo[1])
				}
			}
		} else {
			// Jika tidak ada file logo yang diupload, gunakan nilai lama
			if stores != nil && len(stores.Logo) > 0 {
				fileURL = stores.Logo[0]
				fileID = stores.Logo[1]
				fileNAME = stores.Logo[2]
				projectID = stores.Logo[3]
			}
		}

		slug := utils.CreateSlug(name)
		now := time.Now().Format(time.RFC3339)
		created := now

		//
		table, err := strconv.Atoi(tabl)
		if err != nil {
			http.Redirect(w, r, "/app/store?error=invalid table number", http.StatusSeeOther)
			return
		}

		for i := 1; i <= table; i++ {
			err := client.CheckAndCreateTable(os.Getenv("TABLES"), models.GlobalSessionData.UserId, i)
			if err != nil {
				log.Println("STORE ERROR : " + err.Error())
			}
		}

		updates := models.Store{
			Name:    name,
			Address: []string{city, address},
			Logo:    []string{fileURL, fileID, fileNAME, projectID},
			Slug:    slug,
			Table:   table,

			Merchant: []string{user_id, merchantData[0].MerchantName},
			Owner:    []string{merchantData[0].OwnerId, ownerData.OwnerName},

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

		http.Redirect(w, r, "/app/store?msg=store updated", http.StatusSeeOther)
	}
}

func Billing(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		user_id := models.GlobalSessionData.UserId

		if user_id == "" {
			http.Redirect(w, r, "/app/signout?error=sesi habis", http.StatusSeeOther)
			return
		}

		stores, _ := client.StoreByUserID(os.Getenv("STORES"), user_id)

		if stores == nil {
			stores = &models.Store{} // pakai ini karena limit 1
		}

		packg, _ := client.ListPackage(os.Getenv("PACKAGES")) // ini list

		data := models.PublicData{
			Title: "Billing",
			Data: map[string]interface{}{
				"stores": stores,
				"packg":  packg,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/billing.html", data)
		return
	}
}
