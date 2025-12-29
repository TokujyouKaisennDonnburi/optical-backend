package tool

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CalendarListTool struct {
	agentQueryRepository repository.AgentQueryRepository
	userId               uuid.UUID
	streamFn             func(context.Context, []byte) error
}

func NewCalendarListTool(
	agentQueryRepository repository.AgentQueryRepository,
	userId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*CalendarListTool, error) {
	if agentQueryRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &CalendarListTool{
		agentQueryRepository: agentQueryRepository,
		userId:               userId,
		streamFn:             streamFn,
	}, nil
}

func (t CalendarListTool) Name() string {
	// カレンダー一覧取得ツール
	return "list_calendars"
}

func (t CalendarListTool) Description() string {
	// ユーザーのカレンダーを一覧取得します。
	return "Retrieves a list of all available calendars for the user."
}

func (t CalendarListTool) Strict() bool {
	return false
}

func (t CalendarListTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{},
	}
}

func (t CalendarListTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Debug("calendar list called")
	err := t.streamFn(ctx, statusChunk("fetching_calendars"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	calendars, err := t.agentQueryRepository.FindCalendarByUserId(ctx, t.userId)
	if err != nil {
		return "", err
	}
	output, err := json.Marshal(calendars)
	if err != nil {
		return "", err
	}
	logrus.WithField("len", len(calendars)).Debug("calendar list tool called")
	return string(output), nil
}
