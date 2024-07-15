package services

import (
	"encoding/json"
	"fmt"
	"io"
	"pos/models"
	"time"
)

func (c *AppwriteClient) CheckoutToday(collectionID, userID string) ([]models.Checkout, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	today := time.Now().Format("2006-01-02") // Menggunakan format YYYY-MM-DD
	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"created_at\",\"values\":[\"%s\"]}&queries[1]={\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}", today, userID)
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
		Documents []models.Checkout `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}
