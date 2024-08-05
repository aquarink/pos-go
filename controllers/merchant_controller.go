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
		stores, err := client.ListStores(os.Getenv("STORES"))
		if err != nil {
			http.Redirect(w, r, "/app/dashboard?error=failed to load merchant", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Store",
			Data: map[string]interface{}{
				"stores": stores,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/merchant.html", data)
		return
	}
}
