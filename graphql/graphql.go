package graphql

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Client struct {
	BaseURL   string
	AuthToken string
}

func NewClient(baseURL, authToken string) *Client {
	return &Client{
		BaseURL:   baseURL,
		AuthToken: authToken,
	}
}

type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func (c *Client) Send(request Request, result interface{}) error {
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return err
	}

	log.Printf("Request Payload: %s\n", string(body))

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	log.Printf("Request Headers: %+v\n", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return err
	}

	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response Body: %s\n", string(respBody))

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d, body: %s\n", resp.StatusCode, string(respBody))
		return err
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		log.Printf("Failed to unmarshal response: %v\n", err)
		return err
	}

	log.Printf("Unmarshalled Result: %+v\n", result)
	return nil
}
