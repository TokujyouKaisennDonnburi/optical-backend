package command

import (
	"context"
	"fmt"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/tool"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/openrouter"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	SYSTEM_PROMPT = `
あなたは高度に進化した予定管理エージェントです。

## 回答ルール
- ユーザーの要求や質問に応じて予定の分析をします。
- 現在の日時を考慮して分析を行います。
- 分析や質問の結果を丁寧な口調で説明します。
- 予定が終日の場合は、時間を考慮せず、開始日と終了日のみを考慮します。
- ユーザーが予定の分析に該当しない要求をした場合は、正しい要求の仕方を出力します。

## 現在の日時

%s

`
)

type AgentCommand struct {
	db         *sqlx.DB
	openRouter *openrouter.OpenRouter
}

func NewAgentCommand(
	db *sqlx.DB,
	openRouter *openrouter.OpenRouter,
) *AgentCommand {
	if db == nil {
		panic("db is nil")
	}
	if openRouter == nil {
		panic("openRouter is nil")
	}
	return &AgentCommand{
		db:         db,
		openRouter: openRouter,
	}
}

type AgentCommandExecInput struct {
	UserInput   string
	UserId      uuid.UUID
	StreamingFn func(context.Context, []byte) error
}

func (c *AgentCommand) Exec(ctx context.Context, input AgentCommandExecInput) error {
	// ツール定義
	tools := []openrouter.Tool{
		tool.NewEventSearchTool(c.db, input.UserId),
	}
	systemPrompt := fmt.Sprintf(SYSTEM_PROMPT, time.Now())
	messages := []openrouter.Message{
		openrouter.SystemMessage(systemPrompt),
		openrouter.UserMessage(input.UserInput),
	}
	_, err := c.openRouter.WithTools(tools).ChainStream(ctx, messages, input.StreamingFn)
	return err
}
