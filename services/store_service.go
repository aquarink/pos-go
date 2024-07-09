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

func (c *AppwriteClient) GetStoreByUserID(collectionID, userID string) (*models.Store, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[]=user=%s", userID)
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

func (c *AppwriteClient) GetStoreByID(collectionID, id string) (*models.Store, error) {
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
		"user":    []string{stores.User[0], stores.User[1]},
		"name":    stores.Name,
		"address": []string{stores.Address[0], stores.Address[1]},
		"logo":    []string{stores.Logo[0], stores.Logo[1]},
		"slug":    stores.Slug,
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
	store, err := c.GetStoreByUserID(collectionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find store: %v", err)
	}

	docID := store.ID
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, docID)

	dt := map[string]interface{}{
		"user":    []string{stores.User[0], stores.User[1]},
		"name":    stores.Name,
		"address": []string{stores.Address[0], stores.Address[1]},
		"logo":    []string{stores.Logo[0], stores.Logo[1]},
		"slug":    stores.Slug,
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

func (c *AppwriteClient) UploadLogo(bucketID, fileID string, filePath string) (string, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/files", c.Endpoint, bucketID)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// pake file ID
	fw, err := w.CreateFormField("fileId")
	if err != nil {
		return "", fmt.Errorf("failed to create form field for file ID: %w", err)
	}

	// bikin file data
	fw, err = w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-Appwrite-Project", c.ProjectID)
	req.Header.Set("X-Appwrite-Key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload file, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var fileResponse struct {
		FileID string `json:"$id"`
	}

	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	fileURL := fmt.Sprintf("%s/storage/buckets/%s/files/%s/view", c.Endpoint, bucketID, fileResponse.FileID)
	return fileURL, nil
}
