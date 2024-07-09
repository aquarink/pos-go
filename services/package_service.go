package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) ListPackage(collectionID string) ([]models.Packages, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	req, err := c.newRequest("GET", url, nil)
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
		Documents []models.Packages `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) CreatePackage(collectionID string, category models.Packages) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"name":        category.Name,
		"price":       category.Price,
		"cashier":     category.CashierAvailable,
		"product":     category.ProductAvailable,
		"description": category.Description,
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

	req, err := c.newRequest("POST", url, jsons)
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

func (c *AppwriteClient) PackageById(collectionID, id string) (*models.Packages, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	req, err := c.newRequest("GET", url, nil)
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

	var mdl models.Packages
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) PackageByName(collectionID, name string) (*models.Packages, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	req, err := c.newRequest("GET", url, nil)
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
		Documents []models.Packages `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, doc := range response.Documents {
		if doc.Name == name {
			return &doc, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (c *AppwriteClient) UpdatePackage(collectionID, id string, category models.Packages) (*models.Packages, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	dt := map[string]interface{}{
		"name":        category.Name,
		"price":       category.Price,
		"cashier":     category.CashierAvailable,
		"product":     category.ProductAvailable,
		"description": category.Description,
	}
	updateData := map[string]interface{}{
		"data": dt,
	}

	jsons, err := json.Marshal(updateData)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("PATCH", url, jsons)
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

	var mdl models.Packages
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) DeletePackage(collectionID, id string) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	req, err := c.newRequest("DELETE", url, nil)
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
