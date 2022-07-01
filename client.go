package gpt3client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	apiKey = os.Getenv("OPEN_AI_APIKEY")
	engine = "text-davinci-002"
)

type request struct {
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             float64 `json:"top_p"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
}

type response struct {
	ID      string `json:"id"`
	Choices []struct {
		Text     string    `json:"text"`
		Logprobs []float64 `json:"logprobs"`
		Ranking  float64   `json:"ranking"`
	} `json:"choices"`
}

func SendOpenAIPrompt(prompt string) (outputResponse string, newContent string) {

	request := request{
		Prompt:           prompt,
		Temperature:      0.7,
		MaxTokens:        256,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	requestBody := bytes.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.openai.com/v1/engines/%s/completions", engine), requestBody)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response response
	json.Unmarshal(responseBody, &response)

	outputResponse = prompt
	for _, choice := range response.Choices {
		// fmt.Println(choice.Text)
		if choice.Text == "" {
			continue
		}
		newContent = choice.Text
		outputResponse = fmt.Sprintf("%s%s", outputResponse, newContent)
	}

	return outputResponse, newContent
}

func SendOpenAIPromptStreaming(inputPrompt string, iterations int) (outputResponse string) {
	var newContent string
	var finalContent string
	for i := 0; i < iterations; i++ {

		inputPrompt, newContent = SendOpenAIPrompt(inputPrompt)
		finalContent = fmt.Sprintf("%s%s", finalContent, newContent)
	}

	return finalContent
}
