package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
)

type OpenRouterChatRequest struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
}

type OpenRouterChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (r *OptionAgentOpenRouterRepository) SuggestOptions(ctx context.Context, userPrompt string, options []option.Option) ([]option.Option, error) {
	systemPrompt, err := getOptionSuggestPrompt(options)
	if err != nil {
		return nil, err
	}
	request := OpenRouterChatRequest{
		Model: "openai/gpt-oss-20b:free",
		Messages: []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
	}
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(reqJson),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("openrouter status error: %d", resp.StatusCode)
	}
	var respBody OpenRouterChatResponse
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	var results []OptionAgentSuggestResponse
	err = json.Unmarshal([]byte(respBody.Choices[0].Message.Content), &results)
	if err != nil {
		return nil, err
	}
	options = make([]option.Option, len(results))
	for i, result := range results {
		options[i] = option.Option{
			Id:   result.Id,
			Name: result.Name,
		}
	}
	return options, nil
}
