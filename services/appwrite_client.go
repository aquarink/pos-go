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
)

type AppwriteClient struct {
	Client     *http.Client
	Endpoint   string
	ProjectID  string
	APIKey     string
	DatabaseID string
}

func NewAppwriteClient(endpoint, projectID, apiKey, databaseID string) *AppwriteClient {
	return &AppwriteClient{
		Client:     &http.Client{},
		Endpoint:   endpoint,
		ProjectID:  projectID,
		APIKey:     apiKey,
		DatabaseID: databaseID,
	}
}

func (c *AppwriteClient) kirimRequestKeAppWrite(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Appwrite-Key", c.APIKey)
	req.Header.Add("X-Appwrite-Project", c.ProjectID)
	return req, nil
}

func (c *AppwriteClient) FileUpload(bucketID string, filePath string) (string, string, string, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/files", c.Endpoint, bucketID)

	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Write the file ID
	fw, err := w.CreateFormField("fileId")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create form field for file ID: %w", err)
	}
	_, err = fw.Write([]byte("unique()"))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to write file ID: %w", err)
	}

	// Write the file data
	fw, err = w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to copy file data: %w", err)
	}

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-Appwrite-Project", c.ProjectID)
	req.Header.Set("X-Appwrite-Key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("failed to upload file, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var fileResponse struct {
		FileID   string `json:"$id"`
		FileName string `json:"name"`
	}

	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	fileURL := fmt.Sprintf("%s/storage/buckets/%s/files/%s/view", c.Endpoint, bucketID, fileResponse.FileID)
	return fileURL, fileResponse.FileID, fileResponse.FileName, nil
}

func (c *AppwriteClient) FileRemove(bucketID, fileID string) error {
	url := fmt.Sprintf("%s/storage/buckets/%s/files/%s", c.Endpoint, bucketID, fileID)

	req, err := c.kirimRequestKeAppWrite("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete file, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
