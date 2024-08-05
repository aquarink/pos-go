package services

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pos/models"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/skip2/go-qrcode"
)

func (c *AppwriteClient) ListTables(collectionID, merchantId string) ([]models.Table, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}", merchantId)
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
		Documents []models.Table `json:"documents"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Documents, nil
}

// CheckAndCreateTable mengecek apakah meja sudah ada, jika belum maka membuatnya dan menghasilkan kode QR.
func (c *AppwriteClient) CheckAndCreateTable(collectionID, userID string, tableNo int) error {
	url := fmt.Sprintf("%s/databases/%s/collections/%s/documents", c.Endpoint, c.DatabaseID, collectionID)

	query := fmt.Sprintf("?queries[0]={\"method\":\"equal\",\"attribute\":\"user_id\",\"values\":[\"%s\"]}&queries[1]={\"method\":\"equal\",\"attribute\":\"table_no\",\"values\":[%d]}", userID, tableNo)
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
		// Path relatif ke folder proyek
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get project directory: %v", err)
		}

		// Ensure the tmp directory exists
		tmpDir := filepath.Join(projectDir, "tmp")
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			err = os.Mkdir(tmpDir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create tmp directory: %v", err)
			}
		}

		generateCode := strconv.Itoa(tableNo) + userID
		qrCodeFilePath := filepath.Join(projectDir, "tmp", fmt.Sprintf("%s.png", generateCode))
		// err = GenerateQRCode(generateCode, qrCodeFilePath)
		err = GenerateQRCodeCustom(generateCode, qrCodeFilePath, tableNo)
		if err != nil {
			return err
		}

		qrCodeURL, qrCodeID, qrCodeName, projectID, err := c.FileUpload(os.Getenv("TABLES_BUCKET"), qrCodeFilePath)
		if err != nil {
			log.Println("699999 >> " + err.Error())
			return err
		}

		// Hapus file QR code setelah upload berhasil
		err = os.Remove(qrCodeFilePath)
		if err != nil {
			return err
		}

		now := time.Now().Format(time.RFC3339)
		table := models.Table{
			UserId:    userID,
			TableNo:   tableNo,
			Code:      generateCode,
			CodeImage: []string{qrCodeURL, qrCodeID, qrCodeName, projectID},
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

func GenerateQRCodeCustom(text, filePath string, tableNo int) error {
	// Generate QR code
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	qrImage := qr.Image(256)

	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get project directory: %v", err)
	}

	// Load font
	fontPath := projectDir + "/assets/TT-Neoris-Trial-Bold.ttf"
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("failed to load font: %v", err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse font: %v", err)
	}

	// Create a colored image to draw the QR code
	// totalHeight := 256 + 50
	rgba := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(rgba, rgba.Bounds(), qrImage, image.Point{0, 0}, draw.Src)

	// Add text to the QR code
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(24)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.Black))

	// Calculate the position to center the text
	tulisan := fmt.Sprintf("%d", tableNo)
	textWidth := len(tulisan) * 12
	pt := freetype.Pt((256-textWidth)/2, 250)
	_, err = c.DrawString(tulisan, pt)
	if err != nil {
		return fmt.Errorf("failed to draw string: %v", err)
	}

	// Save the image to a file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, rgba)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}

	return nil
}
