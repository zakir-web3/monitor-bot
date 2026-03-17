package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSendTelegram(t *testing.T) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if token == "" || chatID == "" {
		t.Skip("TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg := "<b>Test Message</b>\n\nThis is a test from <code>monitor-bot</code> at " + time.Now().Format(time.RFC3339)
	if err := sendTelegram(ctx, token, chatID, msg); err != nil {
		t.Fatalf("sendTelegram failed: %v", err)
	}
	t.Log("message sent successfully")
}
