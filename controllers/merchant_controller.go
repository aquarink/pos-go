package controllers

import (
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"

	"github.com/gorilla/sessions"
)

func MerchantList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		merch, err := client.ListMerchants(os.Getenv("STORES"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load merchant", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Merchant",
			Data: map[string]interface{}{
				"merchant": merch,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/merchant.html", data)
		return
	}
}
