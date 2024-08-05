package services

import (
	"encoding/json"
	"fmt"
	"io"
	"pos/models"
)

func (c *AppwriteClient) ListMerchants(collectionID string) ([]models.Store, error) {
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
		Documents []models.Store `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *AppwriteClient) MerchantByMerchantId(collectionID, merchant_id string) ([]models.Merchant, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("{\"method\":\"equal\",\"attribute\":\"merchant_id\",\"values\":[\"%s\"]}", merchant_id)
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
		Documents []models.Merchant `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}
