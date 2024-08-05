package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) ListStores(collectionID string) ([]models.Store, error) {
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
		Documents []models.Store `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) StoreByUserID(collectionID, userID string) (*models.Store, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"merchant\",\"values\":[\"%s\"]}", userID)
	url = url + query

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
		Documents []models.Store `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Documents) == 0 {
		return nil, fmt.Errorf("store not found for user_id: %s", userID)
	}

	return &response.Documents[0], nil
}

func (c *AppwriteClient) ListStoreByOwnerID(collectionID, ownerID string) ([]models.Store, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"owner\",\"values\":[\"%s\"]}", ownerID)
	url = url + query

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
		Documents []models.Store `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) StoreByOwnerID(collectionID, ownerID string) (*models.Store, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"owner\",\"values\":[\"%s\"]}", ownerID)
	url = url + query

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
		Documents []models.Store `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Documents) == 0 {
		return nil, fmt.Errorf("store not found for ownerID: %s", ownerID)
	}

	return &response.Documents[0], nil
}

func (c *AppwriteClient) StoreByID(collectionID, id string) (*models.Store, error) {
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

	var mdl models.Store
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) CreateStore(collectionID string, stores models.Store) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"name":     stores.Name,
		"address":  stores.Address,
		"logo":     stores.Logo,
		"slug":     stores.Slug,
		"table":    stores.Table,
		"owner":    stores.Owner,
		"merchant": stores.Merchant,
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

func (c *AppwriteClient) UpdateStore(collectionID, userID string, stores models.Store) (*models.Store, error) {
	store, err := c.StoreByUserID(collectionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find store: %v", err)
	}

	docID := store.ID
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, docID)

	dt := map[string]interface{}{
		"name":     stores.Name,
		"address":  stores.Address,
		"logo":     stores.Logo,
		"slug":     stores.Slug,
		"table":    stores.Table,
		"owner":    stores.Owner,
		"merchant": stores.Merchant,
	}
	updateData := map[string]interface{}{
		"data": dt,
	}

	jsons, err := json.Marshal(updateData)
	if err != nil {
		return nil, err
	}

	req, err := c.kirimRequestKeAppWrite("PATCH", url, jsons)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mdl models.Store
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}
