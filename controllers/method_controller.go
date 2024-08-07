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

func MethodList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		pays, err := client.ListPayment(os.Getenv("PAYMENTS"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load pays", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Payment Method",
			Data: map[string]interface{}{
				"method": pays,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/payment/method_list.html", data)
		return
	}
}

func MethodAdd(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {

		data := models.PublicData{
			Title:   "Add New  Payment Method",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/payment/method_add.html", data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		typeValue := r.FormValue("type")  // gunakan typeValue untuk menghindari kata kunci 'type'
		percent := r.FormValue("percent") // percent
		denom := r.FormValue("denom")     // denom or decimal
		tax := r.FormValue("tax")         // percent

		if name == "" || typeValue == "" {
			http.Redirect(w, r, "/app/method/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		percentFloat64, err := strconv.ParseFloat(percent, 64) // menggunakan ParseFloat untuk float64
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=invalid percent", http.StatusSeeOther)
			return
		}

		denomFloat64, err := strconv.ParseFloat(denom, 64) // menggunakan ParseFloat untuk float64
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=invalid denom", http.StatusSeeOther)
			return
		}

		taxFloat64, err := strconv.ParseFloat(tax, 64) // menggunakan ParseFloat untuk float64
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=invalid tax", http.StatusSeeOther)
			return
		}

		// Casting float64 ke float32
		percentFloat := float32(percentFloat64)
		denomFloat := float32(denomFloat64)
		taxFloat := float32(taxFloat64)

		// INI UPLOTAN
		file, _, err := r.FormFile("icon")
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=failed to upload icon", http.StatusSeeOther)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp("", "")
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=failed to create temp icon", http.StatusSeeOther)
			return
		}
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=failed to save temp icon", http.StatusSeeOther)
			return
		}

		fileURL, fileID, fileNAME, projectID, err := client.FileUpload(os.Getenv("FILES_BUCKET"), tempFile.Name())
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/app/method/add?error=failed to upload icon to server", http.StatusSeeOther)
			return
		}

		now := time.Now().Format(time.RFC3339)

		paymentMethod := models.Method{
			Name:          name,
			Tipe:          typeValue,
			Icon:          []string{fileURL, fileID, fileNAME, projectID},
			TrxFeePercent: percentFloat,
			TrxFeeDenom:   float32(denomFloat),
			TrxTax:        float32(taxFloat),
			Status:        models.StatusActive,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		err = client.CreatePayment(os.Getenv("PAYMENTS"), paymentMethod)
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=failed to create payment method", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/method/list?msg=payment method created successfully", http.StatusSeeOther)
	}
}

func MethodEdit(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/method/list?error=invalid data", http.StatusSeeOther)
			return
		}

		pays, err := client.PaymentByID(os.Getenv("PAYMENTS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/method/list?error=method not found", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "Update Payment Method",
			Data: map[string]interface{}{
				"method": pays,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/payment/method_edit.html", data)
		return
	}
}

func MethodUpdate(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		typeValue := r.FormValue("type")  // gunakan typeValue untuk menghindari kata kunci 'type'
		percent := r.FormValue("percent") // percent
		denom := r.FormValue("denom")     // denom or decimal
		tax := r.FormValue("tax")         // percent

		if name == "" || typeValue == "" {
			http.Redirect(w, r, "/app/method/list?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		pays, err := client.PaymentByID(os.Getenv("PAYMENTS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/method/list?error=method not found", http.StatusSeeOther)
			return
		}

		percentFloat64, err := strconv.ParseFloat(percent, 64) // menggunakan ParseFloat untuk float64
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=invalid percent", http.StatusSeeOther)
			return
		}

		denomFloat64, err := strconv.ParseFloat(denom, 64) // menggunakan ParseFloat untuk float64
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=invalid denom", http.StatusSeeOther)
			return
		}

		taxFloat64, err := strconv.ParseFloat(tax, 64) // menggunakan ParseFloat untuk float64
		if err != nil {
			http.Redirect(w, r, "/app/method/add?error=invalid tax", http.StatusSeeOther)
			return
		}

		// Casting float64 ke float32
		percentFloat := float32(percentFloat64)
		denomFloat := float32(denomFloat64)
		taxFloat := float32(taxFloat64)

		var fileURL string
		var fileID string
		var fileNAME string
		var projectID string

		file, _, err := r.FormFile("icon")
		log.Println(err.Error())
		if err == nil {
			defer file.Close()

			if pays != nil && len(pays.Icon) > 0 {
				_ = client.FileRemove(os.Getenv("FILES_BUCKET"), pays.Icon[1])
			}

			tempFile, err := os.CreateTemp("", "")
			if err != nil {
				http.Redirect(w, r, "/app/method/add?error=failed to create temp icon", http.StatusSeeOther)
				return
			}
			defer tempFile.Close()
			defer os.Remove(tempFile.Name())

			_, err = io.Copy(tempFile, file)
			if err != nil {
				http.Redirect(w, r, "/app/method/add?error=failed to save temp icon", http.StatusSeeOther)
				return
			}

			fileURL, fileID, fileNAME, projectID, err = client.FileUpload(os.Getenv("FILES_BUCKET"), tempFile.Name())
			if err != nil {
				http.Redirect(w, r, fmt.Sprintf("/app/method/edit/%s?error=failed to upload file", id), http.StatusSeeOther)
				return
			}
		} else {
			fileURL = pays.Icon[0]
			fileID = pays.Icon[1]
			fileNAME = pays.Icon[2]
			projectID = pays.Icon[3]
		}

		now := time.Now().Format(time.RFC3339)

		paymentMethodUpdate := models.Method{
			Name:          name,
			Tipe:          typeValue,
			Icon:          []string{fileURL, fileID, fileNAME, projectID},
			TrxFeePercent: percentFloat,
			TrxFeeDenom:   float32(denomFloat),
			TrxTax:        float32(taxFloat),
			Status:        models.StatusActive,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		_, err = client.UpdatePayment(os.Getenv("PAYMENTS"), id, paymentMethodUpdate)
		if err != nil {
			http.Redirect(w, r, "/app/method/list?error=failed to create payment", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/method/list?msg=payment created successfully", http.StatusSeeOther)
	}
}

func MethodDelete(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/method/list?error=invalid data", http.StatusSeeOther)
			return
		}

		err := client.DeletePayment(os.Getenv("PAYMENTS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/method/list?error=payment not found", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/method/list?msg=payment deleted successfully", http.StatusSeeOther)
	}
}
