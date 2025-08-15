package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Token   string
	ChatID  string
	BaseURL string
}

func NewClient(token string, chatID string) *Client {
	return &Client{
		Token:   token,
		ChatID:  chatID,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

func (c *Client) SendMessage(text string) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	body := map[string]interface{}{
		"chat_id":    c.ChatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s | response: %s", resp.Status, string(respData))
	}

	return nil
}
