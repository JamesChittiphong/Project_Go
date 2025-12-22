package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// SendWebhookNotification posts JSON payload to webhookURL with a short timeout.
func SendWebhookNotification(webhookURL string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil
	}
	return nil
}

// GetEnv returns the environment variable value or fallback.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
