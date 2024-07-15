package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/sessions"
)

func Order(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		product, err := client.ProductByUserID(os.Getenv("PRODUCTS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load merchant", http.StatusSeeOther)
			return
		}

		noAntrian := 1
		checkout, err := client.CheckoutToday(os.Getenv("CHECKOUTS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load checkouts", http.StatusSeeOther)
			return
		}

		if checkout != nil {
			if len(checkout) > 0 {
				noAntrian = len(checkout)
			}
		}

		meja := make([]int, 30)
		for i := 0; i < 30; i++ {
			meja[i] = i + 1
		}

		uniqueID := utils.Uniqid(true)

		data := models.PublicData{
			Title: "Order",
			Data: map[string]interface{}{
				"products": product,
				"meja":     meja,
				"unique":   uniqueID,
				"antrian":  noAntrian,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/order.html", data)
		return
	}
}

func Checkout(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodPost {

		// trxID := r.FormValue("trxID")
		// dineType := r.FormValue("dineType")
		// tableNumber := r.FormValue("tableNumber")

		// productIDs := r.FormValue("productIDs")
		// productNames := r.FormValue("productNames")
		// productPrices := r.FormValue("productPrices")
		// productQtys := r.FormValue("productQtys")
		// productNotes := r.FormValue("productNotes")

		// jumlahItems := r.FormValue("jumlahItems")
		// tax := r.FormValue("tax")
		// taxTotal := r.FormValue("taxTotal")
		// totalPayment := r.FormValue("totalPayment")
		// payMethod := r.FormValue("payMethod")
		// changes := r.FormValue("changes")

		product, err := client.ProductByUserID(os.Getenv("PRODUCTS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load merchant", http.StatusSeeOther)
			return
		}

		uniqueID := utils.Uniqid(true)

		data := models.PublicData{
			Title: "Order",
			Data: map[string]interface{}{
				"products": product,
				"unique":   uniqueID,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/order.html", data)
		return
	}
}
