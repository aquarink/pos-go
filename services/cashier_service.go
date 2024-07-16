package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) ListCashier(collectionID string) ([]models.Cashier, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

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
