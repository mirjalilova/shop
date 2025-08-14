package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Token   string
	ChatID  string
	BaseURL string
}

func NewClient(token, chatID string) *Client {
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API error: %s", resp.Status)
	}
	return nil
}
