package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
	"time"
)

func (c *AppwriteClient) CheckoutToday(collectionID, userID string) ([]models.Checkout, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	today := time.Now().Format("02/01/2006")

	queryDate := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"created_date\",\"values\":[\"%s\"]}", today)
	queryUser := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}", userID)

	url = fmt.Sprintf("%s?queries[]=%s&queries[]=%s", url, queryDate, queryUser)

	req, err := c.kirimRequestKeAppWrite("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Documents []models.Checkout `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) CheckoutListByUserId(collectionID, userID string) ([]models.Checkout, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	queryUser := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}", userID)

	url = fmt.Sprintf("%s?queries[]=%s", url, queryUser)

	req, err := c.kirimRequestKeAppWrite("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Documents []models.Checkout `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) CreateCheckout(collectionID string, checkout models.Checkout) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"user_id":        checkout.UserId,
		"queue":          checkout.Queue,
		"trx_id":         checkout.TrxId,
		"dine_type":      checkout.DineType,
		"table_number":   checkout.TableNumber,
		"items":          checkout.Items,
		"total_item":     checkout.TotalItem,
		"tax":            checkout.Tax,
		"tax_total":      checkout.TaxTotal,
		"total_payment":  checkout.TotalPayment,
		"payment_method": checkout.PaymentMethod,
		"change":         checkout.Change,
		"created_date":   checkout.CreatedDate,
		"created_time":   checkout.CreatedTime,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        dt,
		"permissions": []string{"read(\"any\")"},
	}

	checkoutJSON, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("POST", url, checkoutJSON)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create checkout: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) CreateCheckout2(collectionID string, checkout models.Checkout) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"user_id":        checkout.UserId,
		"queue":          checkout.Queue,
		"trx_id":         checkout.TrxId,
		"dine_type":      checkout.DineType,
		"table_number":   checkout.TableNumber,
		"items":          checkout.Items,
		"total_item":     checkout.TotalItem,
		"tax":            checkout.Tax,
		"tax_total":      checkout.TaxTotal,
		"total_payment":  checkout.TotalPayment,
		"payment_method": checkout.PaymentMethod,
		"change":         checkout.Change,
		"created_date":   checkout.CreatedDate,
		"created_time":   checkout.CreatedTime,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        dt,
		"permissions": []string{"read(\"any\")"},
	}

	jsonData, err := json.Marshal(documentData)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Appwrite-Project", c.ProjectID)
	req.Header.Set("X-Appwrite-Key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create checkout, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
