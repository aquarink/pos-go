package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pos/models"
	"strconv"
	"time"

	"github.com/skip2/go-qrcode"
)

// CheckAndCreateTable mengecek apakah meja sudah ada, jika belum maka membuatnya dan menghasilkan kode QR.
func (c *AppwriteClient) CheckAndCreateTable(collectionID, userID string, tableNo int) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}&queries[1]={\"method\":\"equal\",\"attribute\":\"table_no\",\"values\":[\"%d\"]}", userID, tableNo)
	url = url + query

	req, err := c.kirimRequestKeAppWrite("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response struct {
		Documents []models.Table `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if len(response.Documents) == 0 {
		generateCode := strconv.Itoa(tableNo) + userID
		// Generate QR code image
		qrCodeFilePath := fmt.Sprintf("/tmp/qr_table_%d.png", tableNo)
		err := GenerateQRCode(generateCode, qrCodeFilePath)
		if err != nil {
			return err
		}

		qrCodeURL, qrCodeID, qrCodeName, projectID, err := c.FileUpload(os.Getenv("TABLES_BUCKET"), qrCodeFilePath)
		if err != nil {
			return err
		}

		now := time.Now().Format(time.RFC3339)
		table := models.Table{
			UserId:    userID,
			TableNo:   tableNo,
			Code:      generateCode,
			CodeImage: []string{qrCodeURL, qrCodeID, qrCodeName, projectID}, // []string{fileURL, fileID, fileNAME, projectID},
			CreatedAt: now,
			UpdatedAt: now,
		}
		return c.CreateTable(collectionID, table)
	}

	return nil
}

// CreateTable membuat dokumen meja baru di koleksi yang ditentukan.
func (c *AppwriteClient) CreateTable(collectionID string, table models.Table) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	dt := map[string]interface{}{
		"user_id":    table.UserId,
		"table_no":   table.TableNo,
		"code":       table.Code,
		"code_image": table.CodeImage,
		"created_at": table.CreatedAt,
		"updated_at": table.UpdatedAt,
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

// GenerateQRCode menghasilkan kode QR untuk teks yang diberikan dan menyimpannya sebagai file gambar.
func GenerateQRCode(text, filePath string) error {
	err := qrcode.WriteFile(text, qrcode.Medium, 256, filePath)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}
	return nil
}
