package openrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	openRouterAPIUrl = "https://openrouter.ai/api/v1/chat/completions"
	MinimalReasoning = "minimal"
	LowReasoning     = "low"
	MediumReasoning  = "medium"
	HighReasoning    = "high"
)

var (
	ErrModelNotSet          = errors.New("openrouter model is not set")
	ErrMessageNoContent     = errors.New("message content length is zero")
	ErrResponseNoChoices    = errors.New("message content length is zero")
	ErrChatCompletionFailed = errors.New("completion api request failed")
)

type OpenRouterRequest struct {
	Model          string         `json:"model"`
	Provider       Provider       `json:"provider"`
	Message        []Message      `json:"messages"`
	Temperature    float64        `json:"temperature"`
	Stream         bool           `json:"stream,omitempty"`
	Tools          []ToolRequest  `json:"tools,omitempty"`
	ToolChoice     string         `json:"tool_choice,omitempty"`
	Reasoning      string         `json:"reasoning,omitempty"`
	ResponseFormat map[string]any `json:"response_format,omitempty"`
}

type ToolRequest struct {
	Type     string `json:"type"`
	Function struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Parameters  map[string]any `json:"parameters"`
		Strict      bool           `json:"strict,omitempty"`
	} `json:"function"`
}

type OpenRouterResponse struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Message struct {
	Role       string         `json:"role"`
	Content    string         `json:"content"`
	Reasoning  string         `json:"reasoning,omitempty"`
	Name       string         `json:"name,omitempty"`
	ToolCallId string         `json:"tool_call_id,omitempty"`
	ToolCalls  []FunctionCall `json:"tool_calls,omitempty"`
}

func (r *OpenRouter) Fetch(ctx context.Context, messages []Message) (*OpenRouterResponse, error) {
	if r.model == "" {
		return nil, ErrModelNotSet
	}
	if len(messages) == 0 {
		return nil, ErrMessageNoContent
	}
	sendingMessages := make([]Message, len(messages))
	for i, message := range messages{
		message.Reasoning = ""
		sendingMessages[i] = message
	}
	reqBody := OpenRouterRequest{
		Model:          r.model,
		Provider:       r.provider,
		Stream:         false,
		Message:        sendingMessages,
		Temperature:    r.temperature,
		ResponseFormat: r.responseFormat,
		Tools:          toolsToRequests(r.tools),
		ToolChoice:     r.toolChoice,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	logrus.WithField("reqBody", string(reqBytes)).Info("fetch openrouter request")
	req, err := http.NewRequestWithContext(ctx, "POST", openRouterAPIUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		entry := logrus.WithField("statusCode", resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			entry = entry.WithField("responseBody", body)
		}
		entry.Error("openRouter reqeuest error")
		return nil, ErrChatCompletionFailed
	}
	var respBody OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if len(respBody.Choices) == 0 {
		return nil, ErrResponseNoChoices
	}
	logrus.WithField("respBody", respBody).Info("openrouter fetched")
	return &respBody, nil
}

type openRouterStreamChunk struct {
	Choices []struct {
		Delta        Message `json:"delta"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
}

func (r OpenRouter) Stream(
	ctx context.Context,
	messages []Message,
	streamFn func(context.Context, []byte) error,
) (*OpenRouterResponse, error) {
	if r.model == "" {
		return nil, ErrModelNotSet
	}
	if len(messages) == 0 {
		return nil, ErrMessageNoContent
	}
	sendingMessages := make([]Message, len(messages))
	for i, message := range messages{
		message.Reasoning = ""
		sendingMessages[i] = message
	}
	reqBody := OpenRouterRequest{
		Model:          r.model,
		Provider:       r.provider,
		Stream:         true,
		Message:        sendingMessages,
		Temperature:    r.temperature,
		ResponseFormat: r.responseFormat,
		Tools:          toolsToRequests(r.tools),
		ToolChoice:     r.toolChoice,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	logrus.WithField("reqBody", reqBody).Info("stream openrouter request")
	req, err := http.NewRequestWithContext(ctx, "POST", openRouterAPIUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		entry := logrus.WithField("statusCode", resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			entry = entry.WithField("responseBody", string(body))
		}
		entry.Error("openRouter reqeuest error")
		return nil, ErrChatCompletionFailed
	}
	scanner := bufio.NewScanner(resp.Body)
	deltaMessage := Message{}
	finishReason := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		data = strings.TrimSpace(data)
		if data == "[DONE]" {
			logrus.Debug("received: [DONE]")
			break
		}
		if data == ": OPENROUTER PROCESSING" {
			logrus.Debug("OPENROUTER PROCESSING")
			continue
		}
		var chunk openRouterStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			logrus.WithField("chunk", data).WithError(err).Error("chunk unmarshal error")
			continue
		}
		if len(chunk.Choices) == 0 {
			logrus.Error("choices is nil")
			continue
		}
		delta := chunk.Choices[0].Delta
		deltaMessage.Content += delta.Content
		deltaMessage.Reasoning += delta.Reasoning
		deltaMessage.ToolCallId += delta.ToolCallId
		if chunk.Choices[0].FinishReason != "" {
			finishReason = chunk.Choices[0].FinishReason
		}
		logrus.WithField("delta", delta).Debug("chunk received")
		if len(delta.ToolCalls) > 0 {
			for _, tool := range delta.ToolCalls {
				if len(deltaMessage.ToolCalls) <= tool.Index {
					deltaMessage.ToolCalls = append(
						deltaMessage.ToolCalls,
						make([]FunctionCall, tool.Index-len(deltaMessage.ToolCalls)+1)...,
					)
				}
				deltaMessage.ToolCalls[tool.Index].Id += tool.Id
				deltaMessage.ToolCalls[tool.Index].Type += tool.Type
				deltaMessage.ToolCalls[tool.Index].Function.Name += tool.Function.Name
				deltaMessage.ToolCalls[tool.Index].Function.Arguments += tool.Function.Arguments
			}
		}
		if chunk.Choices[0].FinishReason == reason_tool_calls {
			deltaMessage.Role += delta.Role
			defer logrus.WithField("delta", deltaMessage).Info("tool_calling")
			return &OpenRouterResponse{
				Choices: []Choice{
					{
						FinishReason: reason_tool_calls,
						Message:      deltaMessage,
					},
				},
			}, nil
		}
		streamChunk := fmt.Sprintf("{\"content\":\"%s\"}", delta.Content)
		if err := streamFn(ctx, []byte(streamChunk)); err != nil {
			logrus.WithError(err).Error("chunk streamingFn error")
			continue
		}
	}
	logrus.WithField("delta", deltaMessage).WithField("reason", finishReason).Info("received stream")
	return &OpenRouterResponse{
		Choices: []Choice{
			{
				FinishReason: finishReason,
				Message:      deltaMessage,
			},
		},
	}, nil
}
