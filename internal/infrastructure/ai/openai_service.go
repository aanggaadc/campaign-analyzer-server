package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"campaign-analyzer/internal/usecase"
)

type OpenAIService struct {
	apiKey string
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{apiKey: apiKey}
}

func cleanJSONResponse(content string) string {
	// remove ```json dan ```
	content = strings.TrimSpace(content)

	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	return strings.TrimSpace(content)
}

func (o *OpenAIService) AnalyzeCampaign(input usecase.AnalyzeInput) (usecase.AnalyzeResult, error) {
	// 1. Build prompt
	prompt := BuildAnalyzePrompt(input)

	// 2. Prepare request body
	body := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return usecase.AnalyzeResult{}, err
	}

	// 3. HTTP request
	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return usecase.AnalyzeResult{}, err
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 4. Call API
	client := &http.Client{}
	resp, err := client.Do(req)

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("RAW RESPONSE:", string(bodyBytes))
	fmt.Println("API KEY:", o.apiKey)

	if resp.StatusCode != 200 {
		var errBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errBody)

		return usecase.AnalyzeResult{}, fmt.Errorf("openai error: %v", errBody)
	}

	if err != nil {
		return usecase.AnalyzeResult{}, err
	}
	defer resp.Body.Close()

	var aiResp openAIResponse
	err = json.NewDecoder(resp.Body).Decode(&aiResp)
	if err != nil {
		return usecase.AnalyzeResult{}, err
	}

	if len(aiResp.Choices) == 0 {
		return usecase.AnalyzeResult{}, errors.New("no response from AI")
	}

	content := cleanJSONResponse(aiResp.Choices[0].Message.Content)

	fmt.Println("RAW AI RESPONSE:", content)

	var result usecase.AnalyzeResult
	err = json.Unmarshal([]byte(content), &result)

	if result.Summary == "" {
		return usecase.AnalyzeResult{}, errors.New("empty AI response")
	}

	if err != nil {
		return usecase.AnalyzeResult{}, errors.New("invalid AI JSON response")
	}

	return result, nil
}
