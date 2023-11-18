package discord

import (
	"bytes"
	"gocd/internal/config"
	"net/http"
)

func SendMessage(msg string) error {
	var c = config.GetConfig()

	webhook := c.DiscordWebhook

	if webhook == "" {
		return nil
	}

	data := []byte(`{"content":"` + msg + `"}`)
	body := bytes.NewBuffer(data)

	_, err := http.Post(webhook, "application/json", body)
	return err
}
