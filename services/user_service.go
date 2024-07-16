package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pos/models"
)

func (c *AppwriteClient) CreateUser(collectionID string, user models.User) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	userData := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	documentData := map[string]interface{}{
		"documentId":  "unique()",
		"data":        userData,
		"permissions": []string{"read(\"any\")"},
	}
	userBytes, err := json.Marshal(documentData)
	if err != nil {
		return err
	}

	req, err := c.kirimRequestKeAppWrite("POST", url, userBytes)
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
		return fmt.Errorf("failed to create user: %s", string(body))
	}

	return nil
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

func (c *AppwriteClient) GetUserByEmail(collectionID, email string) ([]models.User, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"email\",\"values\":[\"%s\"]}", email)

	url = fmt.Sprintf("%s?queries[]=%s", url, query)

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
