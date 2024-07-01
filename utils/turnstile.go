package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

const turnstileSiteVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type TurnstileResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes,omitempty"`
}

func VerifyTurnstile(token string) error {
	secretKey := os.Getenv("TURNSTILE_SECRET_KEY")
	if secretKey == "" {
		return errors.New("missing Turnstile secret key")
	}

	resp, err := http.PostForm(turnstileSiteVerifyURL, url.Values{
		"secret":   {secretKey},
		"response": {token},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		return errors.New("invalid Turnstile token")
	}

	return nil
}
