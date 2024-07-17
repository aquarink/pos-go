package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) ListCashier(collectionID, merchantId string) ([]models.Cashier, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"merchant_id\",\"values\":[\"%s\"]}", merchantId)

	url = fmt.Sprintf("%s?queries[]=%s", url, query)

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

	var mdl struct {
		Documents []models.Cashier `json:"documents"`
	}
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return mdl.Documents, nil
}

func (c *AppwriteClient) CreateCashier(collectionID string, cashier models.Cashier) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"merchant_id":   cashier.MerchantId,
		"cashier_id":    cashier.CashierId,
		"cashier_name":  cashier.CashierName,
		"cashier_email": cashier.CashierEmail,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        dt,
		"permissions": []string{"read(\"any\")"},
	}

	jsons, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("POST", url, jsons)
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
		return fmt.Errorf("failed to create: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) CashierById(collectionID, id string) (*models.Cashier, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

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

	var mdl models.Cashier
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) CashierByName(collectionID, name string) ([]models.Cashier, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"cashier_name\",\"values\":[\"%s\"]}", name)

	url = fmt.Sprintf("%s?queries[]=%s", url, query)

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
		Documents []models.Cashier `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) CashierByMerchantId(collectionID, merchantId string) ([]models.Cashier, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"merchant_id\",\"values\":[\"%s\"]}", merchantId)
	url = fmt.Sprintf("%s?queries[]=%s", url, query)

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
		Documents []models.Cashier `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) DeleteCashier(collectionID, id string) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	req, err := c.kirimRequestKeAppWrite("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) UpdateCashierStatus(collectionID, cashierId, status string) error {
	getURL := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)
	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"cashier_id\",\"values\":[\"%s\"]}", cashierId)
	getURL = getURL + query

	req, err := c.kirimRequestKeAppWrite("GET", getURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get cashier: %s", string(body))
	}

	var response struct {
		Documents []struct {
			ID string `json:"$id"`
		} `json:"documents"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if len(response.Documents) == 0 {
		return fmt.Errorf("cashier not found")
	}

	docID := response.Documents[0].ID

	// Update status kasir
	updateURL := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, docID)
	updateData := map[string]interface{}{
		"data": map[string]interface{}{
			"status": status,
		},
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	req, err = c.kirimRequestKeAppWrite("PATCH", updateURL, jsonData)
	if err != nil {
		return err
	}

	resp, err = c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update cashier status: %s", string(body))
	}

	return nil
}
