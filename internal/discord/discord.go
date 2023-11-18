package discord

import (
	"bytes"
	"encoding/json"
	"gocd/internal/config"
	"net/http"
)

type RequestBody struct {
	Content string `json:"content"`
}

func SendMessage(msg string) error {
	var c = config.GetConfig()

	webhook := c.DiscordWebhook

	if webhook == "" {
		return nil
	}

	d := RequestBody{
		Content: msg,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(data)

	_, err = http.Post(webhook, "application/json", body)
	return err
}
