package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
)

const SCHEDULE_ANALYZER_PROMPT = `
あなたは高度に進化した予定分析者です。
あなたは下記の予定を定義したJSON配列の中からユーザーの質問や要求に応じて分析をします。

## 回答ルール
- ユーザーの要求や質問に応じて予定の分析をします。
- ユーザーの要求や質問の結果と関係ない文章は出力しません。
- 分析や質問の結果を丁寧な口調で説明します。
- 予定が終日の場合は、時間を考慮せず、開始日と終了日のみを考慮します。
- ユーザーが予定の分析に該当しない要求をした場合は、正しい要求の仕方を出力します。
- 下記の出力フォーマットに従い出力します。

## 現在の日時

%s

## 予定フォーマット

calendarId : カレンダーID
calendarName : カレンダー名
calendarColor : カレンダーカラー
id : 予定ID
title : 予定名
location : 予定場所
memo : 予定メモ
startAt : 開始日時
endAt : 終了日時
isAllDay : 終日フラグ

## 出力フォーマット

予定の説明文をここに出力します

| 日付 | 時間 | 予定名 | 場所 | メモ |
|------|------|--------|------|------|
| 1月1日 | 09:00 〜 10:00 | 予定名 | 場所 | メモ
| 6月1日 〜 6月6日 | 終日 | 予定名 | 場所 | メモ
| 12月31日 | 終日 | 予定名 | 場所 | メモ
`

func (r *AgentOpenRouterRepository) AnalyzeSchedules(
	ctx context.Context,
	userPrompt string,
	schedules []agent.AnalyzableEvent,
	streamingFn func(context.Context, []byte) error,
) error {
	systemPrompt, err := getAnalyzerPrompt(schedules)
	if err != nil {
		return err
	}
	r.openRouter.WithTemperature(0)
	_, err = r.openRouter.Stream(ctx, []openrouter.Message{
		openrouter.SystemMessage(systemPrompt),
		openrouter.UserMessage(userPrompt),
	}, streamingFn)
	if err != nil {
		return err
	}
	return nil
}

func getAnalyzerPrompt(schedules []agent.AnalyzableEvent) (string, error) {
	json, err := json.Marshal(schedules)
	if err != nil {
		return "", err
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(SCHEDULE_ANALYZER_PROMPT, now) + string(json), nil
}
