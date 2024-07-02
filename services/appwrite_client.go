package services

import (
	"bytes"
	"net/http"
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

func (c *AppwriteClient) newRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Appwrite-Key", c.APIKey)
	req.Header.Add("X-Appwrite-Project", c.ProjectID)
	return req, nil
}
