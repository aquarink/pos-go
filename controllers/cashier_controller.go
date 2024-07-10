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
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/merchant/order.html", data)
		return
	}
}
