package tool

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CalendarDetailTool struct {
	agentQueryRepository repository.AgentQueryRepository
	userId               uuid.UUID
	calendarId           uuid.UUID
	streamFn             func(context.Context, []byte) error
}

func NewCalendarDetailTool(
	agentQueryRepository repository.AgentQueryRepository,
	userId uuid.UUID,
	calendarId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*CalendarDetailTool, error) {
	if agentQueryRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &CalendarDetailTool{
		agentQueryRepository: agentQueryRepository,
		userId:               userId,
		calendarId:           calendarId,
		streamFn:             streamFn,
	}, nil
}

func (t CalendarDetailTool) Name() string {
	// カレンダー一覧取得ツール
	return "detail_calendar"
}

func (t CalendarDetailTool) Description() string {
	// ユーザーのカレンダー詳細を取得します。
	return "Retrieves the detailed information of a specific calendar. Use this to check the calendar's title, memo, location, start_at, end_at, is_allday"
}

func (t CalendarDetailTool) Strict() bool {
	return true
}

func (t CalendarDetailTool) Parameters() map[string]any {
	return map[string]any{
		"type":   "object",
		"properties": map[string]any{},
	}
}

func (t CalendarDetailTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Debug("calendar detail called")
	err := t.streamFn(ctx, statusChunk("fetching_calendars"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	calendar, err := t.agentQueryRepository.FindCalendarByIdAndUserId(ctx, t.userId, t.calendarId)
	if err != nil {
		return "", err
	}
	output, err := json.Marshal(calendar)
	if err != nil {
		return "", err
	}
	logrus.WithField("calendar", calendar).Debug("calendar list tool called")
	return string(output), nil
}
