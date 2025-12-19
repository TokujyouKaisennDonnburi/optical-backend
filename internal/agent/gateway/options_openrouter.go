package gateway

import (
	"context"
	"encoding/json"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"github.com/sirupsen/logrus"
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

func (r *OptionAgentOpenRouterRepository) SuggestOptions(ctx context.Context, userPrompt string, options []option.Option) ([]option.Option, error) {
	systemPrompt, err := getOptionSuggestPrompt(options)
	if err != nil {
		return nil, err
	}
	response, err := r.openRouter.Fetch(ctx, []openrouter.Message{
		openrouter.SystemMessage(systemPrompt),
		openrouter.UserMessage(userPrompt),
	})
	if err != nil {
		return nil, err
	}
	logrus.WithField("content", response.Choices[0].Message.Content).Debug("content fetched")
	var results []OptionAgentSuggestResponse
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &results)
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
