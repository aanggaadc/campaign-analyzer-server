package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"campaign-analyzer/internal/usecase"

	"google.golang.org/genai"
)

type GeminiService struct {
	apiKey string
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{apiKey: apiKey}
}

func (g *GeminiService) AnalyzeCampaign(input usecase.AnalyzeInput) (usecase.AnalyzeResult, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  g.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return usecase.AnalyzeResult{}, err
	}

	prompt := BuildAnalyzePrompt(input)

	contents := []*genai.Content{
		genai.NewContentFromText(prompt, genai.RoleUser),
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		nil,
	)
	if err != nil {
		return usecase.AnalyzeResult{}, err
	}

	if len(resp.Candidates) == 0 {
		return usecase.AnalyzeResult{}, fmt.Errorf("no response from gemini")
	}

	parts := resp.Candidates[0].Content.Parts
	if len(parts) == 0 {
		return usecase.AnalyzeResult{}, fmt.Errorf("empty response parts")
	}

	text := parts[0].Text

	fmt.Println("RAW GEMINI:", text)

	cleaned := cleanJSONResponse(text)

	var result usecase.AnalyzeResult
	err = json.Unmarshal([]byte(cleaned), &result)
	if err != nil {
		return usecase.AnalyzeResult{}, fmt.Errorf("invalid JSON from gemini: %s", cleaned)
	}

	return result, nil
}
