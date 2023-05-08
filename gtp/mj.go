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

type Request struct {
	Action string `json:"action"`
	Prompt string `json:"prompt,omitempty"`
	TaskID string `json:"taskId,omitempty"`
	Index  int    `json:"index,omitempty"`
	State  string `json:"state"`
}

func GetMessageId(prompt string, state string, types string) (string, error) {
	mjImUrl := config.LoadConfig().MjImUrl

	requestData := Request{
		Prompt: prompt,
		Action: types,
		State:  state,
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

	messageId, ok := response["result"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected data format in response: %v", response)
	}
	fmt.Println(messageId)
	return messageId, nil
}

type RequestData2 struct {
	Prompt          string `json:"prompt,omitempty"`
	Content         string `json:"content,omitempty"`
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

	requestData := RequestData2{
		Type:    types,
		State:   state,
		Button:  button,
		TaskId:  taskId,
		Content: taskId + " " + button,
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

	messageId, ok := response["result"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected data format in response: %v", response)
	}
	fmt.Println(messageId)
	return messageId, nil
}
