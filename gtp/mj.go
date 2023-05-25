package gtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/869413421/wechatbot/config"
	"net/http"
)

type requestImg struct {
	Prompt  string `json:"prompt,omitempty"`
	Type    string `json:"type"`
	Webhook string `json:"webhook"`
	State   string `json:"state"`
	MsgHash string `json:"msgHash,omitempty"`
	Index   int64  `json:"index,omitempty"`
	IsAgent bool   `json:"isAgent"`
	Action  string `json:"action"`
}

type ResponseData struct {
	Prompt   string `json:"prompt"`
	PromptEn string `json:"promptEn"`
	TaskId   string `json:"taskId"`
	Length   int64  `json:"length"`
}

func GetMessageId(prompt string, state string, types string) (ResponseData, error) {
	mjImUrl := config.LoadConfig().MjImUrl
	webhook := config.LoadConfig().Webhook

	requestData := requestImg{
		Prompt:  prompt,
		Type:    types,
		Webhook: webhook,
		State:   state,
		IsAgent: true,
		Action:  types,
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

	var responseData ResponseData
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		panic(err)
	}

	fmt.Println(responseData.TaskId)
	return responseData, nil
}

type RequestData2 struct {
	Prompt  string `json:"prompt,omitempty"`
	Type    string `json:"type"`
	Webhook string `json:"webhook"`
	State   string `json:"state"`
	MsgHash string `json:"msgHash,omitempty"`
	Index   int64  `json:"index,omitempty"`
	Button  string `json:"button,omitempty"`
	TaskId  string `json:"taskId,omitempty"`
	IsAgent bool   `json:"isAgent"`
	Action  string `json:"action"`
}

func GetEx(state string, types string, button string, taskId string) (ResponseData, error) {
	mjExImUrl := config.LoadConfig().MjExUrl
	webhook := config.LoadConfig().Webhook

	requestData := RequestData2{
		Type:    types,
		Webhook: webhook,
		State:   state,
		IsAgent: true,
		Button:  button,
		TaskId:  taskId,
		Action:  types,
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

	if err != nil {
		panic(err)
	}
	var responseData ResponseData
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		panic(err)
	}

	fmt.Println(responseData.TaskId)
	return responseData, nil
}
