package example

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ExampleWorkflow shows a minimal client flow for token status and records.
func ExampleWorkflow(apiBase, apiKey, token string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	statusURL := fmt.Sprintf("%s/api/tokens/%s", apiBase, token)
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	recordsURL := fmt.Sprintf("%s/api/tokens/%s/records", apiBase, token)
	req, err = http.NewRequest("GET", recordsURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", apiKey)
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	return nil
}

// ExampleCreateWebhook binds a FIRST_HIT webhook.
func ExampleCreateWebhook(apiBase, apiKey, token, webhookURL, secret string) error {
	body := map[string]string{
		"webhook_url": webhookURL,
		"secret":      secret,
		"mode":        "FIRST_HIT",
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/tokens/%s/webhook", apiBase, token), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
