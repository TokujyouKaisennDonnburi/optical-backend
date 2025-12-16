package gateway

import (
	"context"
	"encoding/json"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"google.golang.org/genai"
)

const OPTION_PROMPT = `
あなたは高度に進化したシステムのオプション提案エージェントです。
あなたは下記のオプションを定義したJSON配列の中からユーザーが求めるオプションを提案します。
指示やガイドラインに従い、必ず次のJSON形式のみで出力してください。

## JSON形式

[
  {
	"id": int,
	"name": string
  }
]

## 提案ルール
- ユーザーが求めているオプションのみを提案する。
- ユーザーが求めているオプションの数を考慮して提案する。
- ユーザーが求めているオプションに該当するものが1つもない場合は何も提案しない。
- 'Description'がユーザーの要求と一致している場合にオプションを提案する。
- Githubに関連するオプションはエンジニア・プログラマーにのみ提案する。

## 提案可能なオプション
`

func getDescriptionMap() map[int32]string {
	return map[int32]string{
		1: "Githubでユーザーがレビューを依頼されているものを確認できます",
		2: "Githubのどのユーザーにレビュー依頼がされているかの付加状況を確認できます",
	}
}

type OptionAgentSuggestResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func (r *OptionAgentGeminiRepository) getConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"id":   {Type: genai.TypeInteger},
					"name": {Type: genai.TypeString},
				},
				Required: []string{"id", "name"},
			},
		},
	}
}

func (r *OptionAgentGeminiRepository) SuggestOptions(ctx context.Context, request string, options []option.Option) ([]option.Option, error) {
	systemPrompt, err := getOptionSuggestPrompt(options)
	if err != nil {
		return nil, err
	}
	parts := []*genai.Part{
		{Text: request},
	}
	config := r.getConfig()
	config.SystemInstruction = &genai.Content{
		Parts: []*genai.Part{
			{Text: systemPrompt},
		},
	}
	resp, err := r.client.Models.GenerateContent(ctx, "gemma-3-4b", []*genai.Content{
		{Parts: parts},
	},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var results []OptionAgentSuggestResponse
	err = json.Unmarshal([]byte(resp.Text()), &results)
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

type OptionSuggestModel struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func getOptionSuggestPrompt(options []option.Option) (string, error) {
	suggests := []OptionSuggestModel{}
	for _, opt := range options {
		desc, ok := getDescriptionMap()[opt.Id]
		if !ok {
			continue
		}
		suggests = append(suggests, OptionSuggestModel{
			Id:          opt.Id,
			Name:        opt.Name,
			Description: desc,
		})
	}
	suggestJson, err := json.Marshal(suggests)
	if err != nil {
		return "", err
	}
	return OPTION_PROMPT + "\n" + string(suggestJson), nil
}
