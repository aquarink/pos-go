package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"pos/models"
)

func (c *AppwriteClient) ListProducts(collectionID string) ([]models.Products, error) {
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

	prodData := map[string]interface{}{
		"name":     product.Name,
		"category": product.Category,
		"price":    product.Price,
		"user_id":  product.UserID,
		"photo":    product.Photo,
		"slug":     product.Slug,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        prodData,
		"permissions": []string{"read(\"any\")"},
	}

	productJSON, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", url, productJSON)
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
		return fmt.Errorf("failed to create product: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) GetProductByID(collectionID, id string) (*models.Products, error) {
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

	var product models.Products
	err = json.Unmarshal(body, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *AppwriteClient) UpdateProduct(collectionID, id string, product models.Products) (*models.Products, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	prodData := map[string]interface{}{
		"name":     product.Name,
		"category": product.Category,
		"price":    product.Price,
		"user_id":  product.UserID,
		"photo":    product.Photo,
		"slug":     product.Slug,
	}
	updateData := map[string]interface{}{
		"data": prodData,
	}

	productJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("PATCH", url, productJSON)
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
		return nil, fmt.Errorf("failed to update product: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var updatedProduct models.Products
	err = json.Unmarshal(body, &updatedProduct)
	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (c *AppwriteClient) DeleteProduct(collectionID, id string) error {
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
		return fmt.Errorf("failed to delete product: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) UploadFile(bucketID, fileID string, filePath string) (string, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/files", c.Endpoint, bucketID)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Write the file ID
	fw, err := w.CreateFormField("fileId")
	if err != nil {
		return "", err
	}
	_, err = fw.Write([]byte(fileID))
	if err != nil {
		return "", err
	}

	// Write the file data
	fw, err = w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return "", err
	}

	// Write the permissions (optional)
	fw, err = w.CreateFormField("permissions[]")
	if err != nil {
		return "", err
	}
	_, err = fw.Write([]byte(`["read(\"any\")"]`))
	if err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-Appwrite-Project", c.ProjectID)
	req.Header.Set("X-Appwrite-Key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload file: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var fileResponse struct {
		FileID string `json:"$id"`
	}
	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return "", err
	}

	return fileResponse.FileID, nil
}
