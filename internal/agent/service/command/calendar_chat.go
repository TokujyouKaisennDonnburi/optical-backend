package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/tool"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	CALENDAR_CHAT_SYSTEM_PROMPT = `
あなたは高度に進化した予定管理エージェントです。
あなたには担当カレンダーがあります。

## 担当カレンダーID

%s

## 予定の管理体制
- ユーザーは複数のカレンダー(calendar)を持っています。
- それぞれのカレンダーに予定(event)があります。

## 回答ルール
- ユーザーの要求や質問に応じてカレンダーや予定の分析や作成をします。
- 予定の分析・作成は日本語で行います。
- 現在の日時を考慮して分析を行います。
- 分析や質問の結果を丁寧な口調で説明します。
- 予定が終日の場合は、時間を考慮せず、開始日と終了日のみを考慮します。
- ユーザーが予定管理と関係しない要求をした場合は、あなたができることを説明します。

## 現在の日時

%s

`
)

type AgentCommandCalendarChatInput struct {
	UserInput   string
	UserId      uuid.UUID
	CalendarId  uuid.UUID
	StreamingFn func(context.Context, []byte) error
}

func (c *AgentCommand) CalendarChat(ctx context.Context, input AgentCommandCalendarChatInput) error {
	userPrompt := strings.TrimSpace(input.UserInput)
	if userPrompt == "" {
		return apperr.ValidationError("invalid user message")
	}
	calendarEventSearchTool, err := tool.NewCalendarEventSearchTool(c.agentQueryRepository, input.UserId, input.CalendarId, input.StreamingFn)
	if err != nil {
		return err
	}
	calendarDetailTool, err := tool.NewCalendarDetailTool(c.agentQueryRepository, input.UserId, input.CalendarId, input.StreamingFn)
	if err != nil {
		return err
	}
	eventCreateTool, err := tool.NewEventCreateTool(c.agentCommandRepository, input.UserId, input.StreamingFn)
	if err != nil {
		return err
	}
	optionListTool, err := tool.NewOptionListTool(c.agentQueryRepository, input.UserId, input.StreamingFn)
	if err != nil {
		return err
	}
	calendarOptionListTool, err := tool.NewCalendarOptionListTool(c.agentQueryRepository, input.UserId, input.CalendarId, input.StreamingFn)
	if err != nil {
		return err
	}
	calendarOptionUpdateTool, err := tool.NewCalendarOptionUpdateTool(c.agentCommandRepository, input.UserId, input.CalendarId, input.StreamingFn)
	if err != nil {
		return err
	}
	// ツール定義
	tools := []openrouter.Tool{
		calendarEventSearchTool,
		calendarDetailTool,
		eventCreateTool,
		optionListTool,
		calendarOptionListTool,
		calendarOptionUpdateTool,
	}
	systemPrompt := fmt.Sprintf(CALENDAR_CHAT_SYSTEM_PROMPT, input.CalendarId.String(), time.Now())
	messages := []openrouter.Message{
		openrouter.SystemMessage(systemPrompt),
		openrouter.UserMessage(userPrompt),
	}
	return c.transactor.Transact(ctx, func(ctx context.Context) error {
		// 生成中ステータスを送信
		err = input.StreamingFn(ctx, []byte("{\"status\": \"analyzing\"}"))
		if err != nil {
			logrus.WithError(err).Error("status streaming error")
		}
		_, err = c.openRouter.WithTools(tools).ChainStream(ctx, messages, input.StreamingFn)
		return err
	})
}
