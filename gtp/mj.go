package gtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/869413421/wechatbot/config"
	"net/http"
)

type requestImg struct {
	Prompt          string `json:"prompt,omitempty"`
	Type            string `json:"type"`
	WebhookOverride string `json:"webhookOverride"`
	State           string `json:"state"`
	MsgHash         string `json:"msgHash,omitempty"`
	Index           int64  `json:"index,omitempty"`
	IsAgent         bool   `json:"isAgent"`
}

func GetMessageId(prompt string, state string, types string) (string, error) {
	mjImUrl := config.LoadConfig().MjImUrl
	webhook := config.LoadConfig().Webhook

	requestData := requestImg{
		Prompt:          prompt,
		Type:            types,
		WebhookOverride: webhook,
		State:           state,
		IsAgent:         true,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", mjImUrl, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	messageId, ok := response["messageId"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected data format in response: %v", response)
	}
	fmt.Println(messageId)
	return messageId, nil
}

type RequestData2 struct {
	Prompt          string `json:"prompt,omitempty"`
	Type            string `json:"type"`
	WebhookOverride string `json:"webhookOverride"`
	State           string `json:"state"`
	MsgHash         string `json:"msgHash,omitempty"`
	Index           int64  `json:"index,omitempty"`
	Button          string `json:"button,omitempty"`
	TaskId          string `json:"taskId,omitempty"`
	IsAgent         bool   `json:"isAgent"`
}

func GetEx(state string, types string, button string, taskId string) (string, error) {
	mjExImUrl := config.LoadConfig().MjExUrl
	webhook := config.LoadConfig().Webhook

	requestData := RequestData2{
		Type:            types,
		WebhookOverride: webhook,
		State:           state,
		IsAgent:         true,
		Button:          button,
		TaskId:          taskId,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", mjExImUrl, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	messageId, ok := response["messageId"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected data format in response: %v", response)
	}
	fmt.Println(messageId)
	return messageId, nil
}
