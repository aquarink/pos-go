package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) ListProducts(collectionID string) ([]models.Products, error) {
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
		Documents []models.Products `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) CreateProduct(collectionID string, product models.Products) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"name":     product.Name,
		"category": product.Category,
		"price":    product.Price,
		"user_id":  product.UserID,
		"photo":    product.Photo,
		"slug":     product.Slug,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        dt,
		"permissions": []string{"read(\"any\")"},
	}

	productJSON, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("POST", url, productJSON)
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

func (c *AppwriteClient) GetProductByID(collectionID, id string) (*models.Products, error) {
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

	var mdl models.Products
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) UpdateProduct(collectionID, id string, product models.Products) (*models.Products, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	dt := map[string]interface{}{
		"name":     product.Name,
		"category": product.Category,
		"price":    product.Price,
		"user_id":  product.UserID,
		"photo":    product.Photo,
		"slug":     product.Slug,
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

	var mdl models.Products
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) DeleteProduct(collectionID, id string) error {
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

func (c *AppwriteClient) ProductByUserID(collectionID, userID string) ([]models.Products, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}", userID)
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
		Documents []models.Products `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) ProductByName(collectionID, name, id string) ([]models.Products, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	queryName := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"name\",\"values\":[\"%s\"]}", name)
	queryUser := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}", id)

	url = fmt.Sprintf("%s?queries[]=%s&queries[]=%s", url, queryName, queryUser)

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
		Documents []models.Products `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}
