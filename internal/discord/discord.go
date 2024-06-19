package discord

import (
	"bytes"
	"encoding/json"
	"gocd/internal/config"
	"io"
	"net/http"
)

type ResponseBody struct {
	ID string `json:"id"`
}

type RequestBody struct {
	Content string `json:"content"`
}

func SendMessage(msg string) (string, error) {
	var c = config.GetConfig()

	webhook := c.DiscordWebhook

	if webhook == "" {
		return "", nil
	}

	d := RequestBody{
		Content: msg,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	resp, err := http.Post(webhook+"?wait=true", "application/json", body)
	if err != nil {
		return "", err
	}
	respData, _ := io.ReadAll(resp.Body)

	var r ResponseBody

	json.Unmarshal(respData, &r)

	return r.ID, err
}

func UpdateMessage(id, msg string) (string, error) {
	var c = config.GetConfig()

	webhook := c.DiscordWebhook

	if webhook == "" {
		return "", nil
	}

	d := RequestBody{
		Content: msg,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("PATCH", webhook+"/messages/"+id, body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respData, _ := io.ReadAll(resp.Body)

	var r ResponseBody

	json.Unmarshal(respData, &r)

	return r.ID, err
}
