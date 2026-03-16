package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendTelegram(token, chatID, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	body, err := json.Marshal(map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status %d", resp.StatusCode)
	}
	return nil
}
