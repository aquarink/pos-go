package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/sessions"
)

func TransactionList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		checkout, err := client.CheckoutListByRoleByUserId(os.Getenv("CHECKOUTS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/transaction?error=failed to load transaction", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Transaction",
			Data: map[string]interface{}{
				"checkout": checkout,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/transaction.html", data)
		return
	}
}

// OWNER

func TransactionForOwnerList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		checkout, err := client.CheckoutListByRoleByUserId(os.Getenv("CHECKOUTS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/owner/transaction?error=failed to load transaction", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Transaction",
			Data: map[string]interface{}{
				"checkout": checkout,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/owner/transaction.html", data)
		return
	}
}
