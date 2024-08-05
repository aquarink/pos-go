package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pos/models"
	"time"

	"github.com/google/uuid"
)

type IPaymuClient struct {
	Endpoint string
	VA       string
	APIKey   string
}

func NewIPaymuClient() *IPaymuClient {
	return &IPaymuClient{
		Endpoint: os.Getenv("IPAYMU_PAY_URL"),
		VA:       os.Getenv("IPAYMU_VA"),
		APIKey:   os.Getenv("IPAYMU_KEY"),
	}
}

func getSHA256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateHMACSHA256(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func (ipaymu *IPaymuClient) generateSignature(method string, requestBody []byte) (string, error) {
	hashedBody := getSHA256Hash(string(requestBody))
	stringToSign := method + ":" + ipaymu.VA + ":" + hashedBody + ":" + ipaymu.APIKey
	signature := generateHMACSHA256(ipaymu.APIKey, stringToSign)
	return signature, nil
}

func (ipaymu *IPaymuClient) ListPaymentChannels() (map[string]interface{}, error) {
	url := ipaymu.Endpoint + "-channels"
	timestamp := time.Now().Format("20060102150405")

	// Create empty body for GET request
	requestBody := []byte("{}")

	signature, err := ipaymu.generateSignature("GET", requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("va", ipaymu.VA)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("signature", signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response, nil
}

func (ipaymu *IPaymuClient) RedirectPayment(totalPrice int, packg models.Packages, owner models.Owner) (map[string]interface{}, error) {
	url := ipaymu.Endpoint
	timestamp := time.Now().Format("20060102150405")
	generateId := uuid.New().String()

	paymentData := map[string]interface{}{
		"product":       []string{"Paket " + packg.Name},
		"qty":           []int{1},
		"price":         []int{totalPrice},
		"description":   []string{packg.Description},
		"returnUrl":     "https://your-website.com/return/" + generateId,
		"notifyUrl":     "https://your-website.com/return/" + generateId,
		"cancelUrl":     "https://your-website.com/return/" + generateId,
		"referenceId":   generateId,
		"buyerName":     owner.OwnerName,
		"buyerEmail":    owner.OwnerId + "@pos.com",
		"buyerPhone":    "081234567890",
		"paymentMethod": "mpm",
	}

	jsonData, err := json.Marshal(paymentData)
	if err != nil {
		return nil, err
	}
	signature, err := ipaymu.generateSignature("POST", jsonData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("signature", signature)
	req.Header.Add("va", ipaymu.VA)
	req.Header.Add("timestamp", timestamp)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
