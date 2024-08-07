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

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func Order(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		var products []models.Products
		var stores []models.Store

		cashiers, err := client.CashierByCashierId(os.Getenv("CASHIERS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/signout?error=failed to load cashier id", http.StatusSeeOther)
			return
		}

		if len(cashiers) > 0 {
			for _, cashier := range cashiers {
				cashierProducts, err := client.ProductByUserID(os.Getenv("PRODUCTS"), cashier.MerchantId)
				if err != nil {
					http.Redirect(w, r, "/app/signout?error=failed to load merchant", http.StatusSeeOther)
					return
				}
				products = append(products, cashierProducts...)

				//

				storeData, err := client.StoreByUserID(os.Getenv("STORES"), cashier.MerchantId)
				if err != nil {
					http.Redirect(w, r, "/app/signout?error=failed to load store id", http.StatusSeeOther)
					return
				}
				stores = append(stores, *storeData)
			}
		} else {
			http.Redirect(w, r, "/app/signout?error=failed to load cashier", http.StatusSeeOther)
			return
		}

		noAntrian := 1
		checkout, err := client.CheckoutToday(os.Getenv("CHECKOUTS"), models.GlobalSessionData.UserId)
		if err != nil {
			log.Println("CASHIER ERROR : " + err.Error())
		}

		if checkout != nil {
			if len(checkout) > 0 {
				noAntrian = len(checkout) + 1
			}
		}

		meja := make([]int, stores[0].Table)
		for i := 0; i < stores[0].Table; i++ {
			meja[i] = i + 1
		}

		uniqueID := utils.Uniqid(true)

		data := models.PublicData{
			Title: "Order",
			Data: map[string]interface{}{
				"products": products,
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

		// CARI
		// cashier by models.GlobalSessionData.UserId ambil CashierName
		cashierData, err := client.CashierByCashierId(os.Getenv("CASHIERS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed cashir data", http.StatusSeeOther)
			return
		}

		// merchant di merchants by cashier. ambil MerchantName
		merchantData, err := client.MerchantByMerchantId(os.Getenv("MERCHANTS"), cashierData[0].MerchantId)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed merchant data", http.StatusSeeOther)
			return
		}

		// owner di users by merchant.OwnerId ambil name
		ownerData, err := client.UserByID(os.Getenv("USERS"), merchantData[0].OwnerId)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed owner data", http.StatusSeeOther)
			return
		}

		checkout := models.Checkout{
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

			CashierData:  []string{cashierData[0].CashierId, cashierData[0].CashierName},
			MerchantData: []string{cashierData[0].MerchantId, merchantData[0].MerchantName},
			OwnerData:    []string{merchantData[0].OwnerId, ownerData.Name},

			CreatedDate: time.Now().Format("02/01/2006"),
			CreatedTime: time.Now().Format("15:04:05"),
		}

		err = client.CreateCheckout2(os.Getenv("CHECKOUTS"), checkout)
		if err != nil {
			http.Redirect(w, r, "/app/order?error=failed to create checkout", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/order?msg=checkout created successfully", http.StatusSeeOther)
	}
}

//

func CashierList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		if models.GlobalSessionData.Role == "Merchant" || models.GlobalSessionData.Role == "merchant" {
			cashier, err := client.ListCashier(os.Getenv("CASHIERS"), models.GlobalSessionData.UserId)
			if err != nil {
				http.Redirect(w, r, "/app/cashier/list?error=failed to load package", http.StatusSeeOther)
				return
			}

			data := models.PublicData{
				Title:   "List of Cashier",
				Data:    map[string]interface{}{"cashiers": cashier},
				Error:   r.URL.Query().Get("error"),
				Msg:     r.URL.Query().Get("msg"),
				Session: models.GlobalSessionData,
			}

			utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/cashier/cashier_list.html", data)
			return
		}

		http.Redirect(w, r, "/app/signout?error=your not allowed", http.StatusSeeOther)
		return
	}
}

func CashierAdd(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {

		data := models.PublicData{
			Title:   "Add New Cashier",
			Data:    map[string]interface{}{},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/cashier/cashier_add.html", data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		repassword := r.FormValue("repassword")

		user_id := models.GlobalSessionData.UserId

		if name == "" || email == "" || password == "" || repassword == "" {
			http.Redirect(w, r, "/app/cashier/add?error=form tidak lengkap", http.StatusSeeOther)
			return
		}

		// PREVENT
		// Merchant data
		merchantData, err := client.MerchantByMerchantId(os.Getenv("MERCHANTS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=failed to load merchant data", http.StatusSeeOther)
			return
		}

		// Owner data
		ownerData, err := client.OwnerDataByOwnerId(os.Getenv("OWNERS"), merchantData[0].OwnerId)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=failed to load owner data", http.StatusSeeOther)
			return
		}

		cashierData, err := client.CashierByMerchantId(os.Getenv("CASHIERS"), user_id)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=data cashier invalid", http.StatusSeeOther)
			return
		}

		maxCashier := ownerData.CashierAvailable

		if len(cashierData) >= maxCashier {
			http.Redirect(w, r, "/app/cashier/add?error=anda tidak dapat menambah kasir, harap upgrade paket", http.StatusSeeOther)
			return
		}

		if len(password) < 8 {
			http.Redirect(w, r, "/app/cashier/add?error=password less than 8 character", http.StatusSeeOther)
			return
		}

		if password != repassword {
			http.Redirect(w, r, "/app/cashier/add?error=password not match", http.StatusSeeOther)
			return
		}

		existingUser, err := client.UserByEmail(os.Getenv("USERS"), email)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=internal server error", http.StatusSeeOther)
			return
		}

		if existingUser != nil {
			http.Redirect(w, r, "/app/cashier/add?error=email sudah ada", http.StatusSeeOther)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=internal server failed", http.StatusSeeOther)
			return
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashedPassword),
			Role:     models.RoleCashier,
		}

		userID, err := client.CreateUser(os.Getenv("USERS"), user)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=internal server error "+err.Error(), http.StatusSeeOther)
			return
		}

		// COLLECTION KASIR
		kasir := models.Cashier{
			MerchantId:   models.GlobalSessionData.UserId,
			CashierId:    userID,
			CashierName:  name,
			CashierEmail: email,
			Status:       models.StatusActive,
		}

		err = client.CreateCashier(os.Getenv("CASHIERS"), kasir)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=internal server error cashier"+err.Error(), http.StatusSeeOther)
			return
		}

		// KIRIM EMAIL
		subject := "Email Verification"
		text := fmt.Sprintf("Hi %s,\n\nThank you for registering with us.", name)
		html := fmt.Sprintf("Hi %s,<br><br>Thank you for registering with us.<br>Click <a href='%s%s'>here</a> to verify your email.", name, os.Getenv("EMAIL_VERIFY_URL"), userID)

		err = utils.SendEmail(email, subject, text, html)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=gagal mengirim email verifikasi", http.StatusSeeOther)
			return
		}

		// MODEL MAILS
		emailDoc := models.Mails{
			UserID:  userID,
			Email:   email,
			Subject: subject,
			Text:    text,
			HTML:    html,
		}

		err = client.CreateEmail(os.Getenv("MAILS"), emailDoc)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/add?error=internal server fails", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/cashier/list?msg=silahkan cek email anda untuk verifikasi", http.StatusSeeOther)
	}
}

func CashierStatus(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/cashier/list?error=invalid data", http.StatusSeeOther)
			return
		}

		cashierData, err := client.CashierById(os.Getenv("CASHIERS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/list?error=gagal mendapatkan data kasir", http.StatusSeeOther)
			return
		}

		stat := models.StatusDeactive
		if cashierData.Status == models.StatusDeactive {
			stat = models.StatusActive
		} else if cashierData.Status == models.StatusActive {
			stat = models.StatusDeactive
		}

		err = client.UpdateCashierStatus(os.Getenv("CASHIERS"), cashierData.CashierId, stat)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/list?error=gagal update status", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/cashier/list?msg=berhasil update status menjadi "+stat, http.StatusSeeOther)
	}
}

func CashierDelete(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Redirect(w, r, "/app/cashier/list?error=invalid data", http.StatusSeeOther)
			return
		}

		err := client.DeleteCashier(os.Getenv("CASHIERS"), id)
		if err != nil {
			http.Redirect(w, r, "/app/cashier/list?error=paket tidak ditemukan", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/app/cashier/list?msg=berhasil menghapus paket", http.StatusSeeOther)
	}
}

// OWNER

func CashierForOwnerList(w http.ResponseWriter, r *http.Request, client *services.AppwriteClient, store *sessions.CookieStore) {
	if r.Method == http.MethodGet {
		cashiers, err := client.CashierByMerchantId(os.Getenv("CASHIERS"), models.GlobalSessionData.UserId)
		if err != nil {
			http.Redirect(w, r, "/app/owner/store?error=failed to load cashier", http.StatusSeeOther)
			return
		}

		data := models.PublicData{
			Title: "List of Cashier",
			Data: map[string]interface{}{
				"cashiers": cashiers,
			},
			Error:   r.URL.Query().Get("error"),
			Msg:     r.URL.Query().Get("msg"),
			Session: models.GlobalSessionData,
		}

		utils.RenderTemplateWithSidebar(w, r, "views/templates/backend.html", "views/pages/owner/cashier.html", data)
		return
	}
}
