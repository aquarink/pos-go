package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"pos/models"
)

func GetAllUsers() ([]models.User, error) {
	collectionID := os.Getenv("USERS_COLLECTION_ID")
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", appwriteEndpoint, appwriteProjectID, collectionID)

	req, err := newAppwriteRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := appwriteClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Documents []models.User `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func CreateUser(user models.User) error {
	collectionID := os.Getenv("USERS_COLLECTION_ID")
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", appwriteEndpoint, appwriteProjectID, collectionID)

	userData := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	userBytes, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	req, err := newAppwriteRequest("POST", url, userBytes)
	if err != nil {
		return err
	}

	resp, err := appwriteClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user: %s", string(body))
	}

	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	collectionID := os.Getenv("USERS_COLLECTION_ID")
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", appwriteEndpoint, appwriteProjectID, collectionID)

	req, err := newAppwriteRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := appwriteClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Documents []models.User `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, doc := range response.Documents {
		if doc.Email == email {
			return &doc, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}
