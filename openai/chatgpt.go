package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ijijni/wechat-gpt/config"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	URL = "https://api.openai.com/v1/chat/completions"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []Choices              `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type Choices struct {
	Index        int
	Message      Msg
	FinishReason string
}

type ChatGPTErrorBody struct {
	Error map[string]interface{} `json:"error"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Messages         []Msg   `json:"messages"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completions sendMsg
func Completions(msg string) (*string, error) {
	apiKey := config.GetOpenAiApiKey()
	if apiKey == nil {
		return nil, errors.New("未配置apiKey")
	}

	requestMsg := make([]Msg, 0, 2)
	system := Msg{
		Role:    "system",
		Content: "你是一个调皮的助手",
	}
	requestMsg = append(requestMsg, system)
	userMsg := Msg{
		Role:    "user",
		Content: msg,
	}
	requestMsg = append(requestMsg, userMsg)

	requestBody := ChatGPTRequestBody{
		Model:            "gpt-3.5-turbo",
		Messages:         requestMsg,
		MaxTokens:        4000,
		Temperature:      0.8,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("request openai json string : %v", string(requestData))

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *apiKey))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = reply + v.Message.Content
			break
		}
	}

	gptErrorBody := &ChatGPTErrorBody{}
	err = json.Unmarshal(body, gptErrorBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(reply) == 0 {
		reply = gptErrorBody.Error["message"].(string)
	}

	log.Printf("gpt response full text: %s \n", reply)
	result := strings.TrimSpace(reply)
	return &result, nil
}
