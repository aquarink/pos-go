package services

import (
	"fmt"
	"net/http"
)

var appwriteClient *http.Client
var appwriteEndpoint string
var appwriteProjectID string
var appwriteAPIKey string

func InitAppwriteClient(endpoint, projectID, apiKey string) {
	appwriteClient = &http.Client{}
	appwriteEndpoint = endpoint
	appwriteProjectID = projectID
	appwriteAPIKey = apiKey
}

func getAppwriteURL(databaseID, collectionID, documentID string) string {
	return fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", appwriteEndpoint, databaseID, collectionID, documentID)
}

func newAppwriteRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Appwrite-Key", appwriteAPIKey)
	req.Header.Add("X-Appwrite-Project", appwriteProjectID)
	return req, nil
}
