package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pos/models"
	"pos/services"
	"pos/utils"
	"strconv"
	"strings"
	"time"

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
			log.Println(" >>>> " + err.Error())
		}

		log.Println(checkout)

		if checkout != nil {
			if len(checkout) > 0 {
				noAntrian = len(checkout) + 1
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
		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed  parse form", http.StatusSeeOther)
			return
		}

		trxID := r.FormValue("trxID")
		dineType := r.FormValue("dineType")
		tableNumber := r.FormValue("tableNumber")

		productIDs := strings.Split(r.FormValue("productIDs"), ",")
		productNames := strings.Split(r.FormValue("productNames"), ",")
		productPrices := strings.Split(r.FormValue("productPrices"), ",")
		productQtys := strings.Split(r.FormValue("productQtys"), ",")
		productNotes := strings.Split(r.FormValue("productNotes"), ",")

		jumlahItems, err := strconv.Atoi(r.FormValue("jumlahItems"))
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed jumlahItems", http.StatusSeeOther)
			return
		}

		tax, err := strconv.ParseFloat(r.FormValue("tax"), 64)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed tax", http.StatusSeeOther)
			return
		}

		taxTotal, err := strconv.ParseFloat(r.FormValue("taxTotal"), 64)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed taxTotal", http.StatusSeeOther)
			return
		}

		totalPayment, err := strconv.ParseFloat(r.FormValue("totalPayment"), 64)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed totalPayment", http.StatusSeeOther)
			return
		}

		payMethod := r.FormValue("payMethod")

		changes, err := strconv.ParseFloat(r.FormValue("changes"), 64)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed changes", http.StatusSeeOther)
			return
		}

		antrian, err := strconv.Atoi(r.FormValue("jumlahItems"))
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed antrian", http.StatusSeeOther)
			return
		}

		items := make([]string, len(productIDs))
		for i := range productIDs {
			items[i] = fmt.Sprintf("%s|%s|%s|%s|%s", productIDs[i], productNames[i], productPrices[i], productQtys[i], productNotes[i])
		}

		checkout := models.Checkout{
			UserId:        models.GlobalSessionData.UserId,
			Queue:         antrian,
			TrxId:         trxID,
			DineType:      dineType,
			TableNumber:   tableNumber,
			Items:         items,
			TotalItem:     jumlahItems,
			Tax:           tax,
			TaxTotal:      taxTotal,
			TotalPayment:  totalPayment,
			PaymentMethod: payMethod,
			Change:        changes,
			CreatedDate:   time.Now().Format("02/01/2006"),
			CreatedTime:   time.Now().Format("15:04:05"),
		}

		err = client.CreateCheckout2(os.Getenv("CHECKOUTS"), checkout)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed to create checkout", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/order?msg=checkout created successfully", http.StatusSeeOther)
	}
}
