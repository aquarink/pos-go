package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FileUploadService struct {
	Client *http.Client
}

func (f *FileUploadService) UploadFile(bucket string, formFile string, r *http.Request) (string, string, string, string, error) {
	file, _, err := r.FormFile(formFile)
	if err != nil {
		return "", "", "", "", err
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", "", "", "", err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", "", "", "", err
	}

	fileURL, fileID, fileName, err := f.uploadToBucket(bucket, tempFile.Name())
	if err != nil {
		return "", "", "", "", err
	}

	projectID := os.Getenv("APPWRITE_PROJECT_ID")

	return fileURL, fileID, fileName, projectID, nil
}

func (f *FileUploadService) uploadToBucket(bucket string, filePath string) (string, string, string, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/files", os.Getenv("APPWRITE_ENDPOINT"), bucket)

	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", "", "", err
	}

	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
	req.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_API_KEY"))

	resp, err := f.Client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("failed to upload file: %s", string(body))
	}

	var response struct {
		FileID   string `json:"$id"`
		FileName string `json:"name"`
		FileURL  string `json:"url"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", "", "", err
	}

	return response.FileURL, response.FileID, response.FileName, nil
}
