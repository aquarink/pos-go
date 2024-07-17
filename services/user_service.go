package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) CreateUser(collectionID string, user models.User) (string, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	userData := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"role":     user.Role,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        userData,
		"permissions": []string{"read(\"any\")"},
	}
	userBytes, err := json.Marshal(documentData)
	if err != nil {
		return "", err
	}

	req, err := c.kirimRequestKeAppWrite("POST", url, userBytes)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create user: %s", string(body))
	}

	var response struct {
		ID string `json:"$id"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

func (c *AppwriteClient) GetAllUsers(collectionID string) ([]models.User, error) {
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
		Documents []models.User `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) GetUserByEmail(collectionID, email string) (*models.User, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"email\",\"values\":[\"%s\"]}", email)
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user by email, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mdl struct {
		Documents []models.User `json:"documents"`
	}
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	if len(mdl.Documents) == 0 {
		return nil, nil
	}

	return &mdl.Documents[0], nil
}

func (c *AppwriteClient) GetUserByID(collectionID, id string) (*models.User, error) {
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

	var mdl models.User
	err = json.Unmarshal(body, &mdl)
	if err != nil {
		return nil, err
	}

	return &mdl, nil
}

func (c *AppwriteClient) UpdateUser(collectionID, id string, user models.User) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, id)

	userData := map[string]interface{}{
		"name":           user.Name,
		"email":          user.Email,
		"password":       user.Password,
		"email_verified": user.EmailVerified,
		"role":           user.Role,
	}
	documentData := map[string]interface{}{
		"data": userData,
	}

	userBytes, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("PATCH", url, userBytes)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update user: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) CreateEmail(collectionID string, email models.Mails) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	emailData := map[string]interface{}{
		"user_id": email.ID,
		"email":   email.Email,
		"subject": email.Subject,
		"text":    email.Text,
		"html":    email.HTML,
	}

	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        emailData,
		"permissions": []string{"read(\"any\")"},
	}

	emailBytes, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("POST", url, emailBytes)
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
		return fmt.Errorf("failed to create email document: %s", string(body))
	}

	return nil
}

func (c *AppwriteClient) VerifyUserEmail(collectionID, userID string) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents/%s", c.Endpoint, c.DatabaseID, collectionID, userID)

	updateData := map[string]interface{}{
		"data": map[string]interface{}{
			"email_verified": true,
		},
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("PATCH", url, jsonData)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to verify email: %s", string(body))
	}

	return nil
}
