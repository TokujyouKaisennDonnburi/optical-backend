package tool

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CalendarEventSearchTool struct {
	agentQueryRepository repository.AgentQueryRepository
	userId               uuid.UUID
	calendarId           uuid.UUID
	streamFn             func(context.Context, []byte) error
}

func NewCalendarEventSearchTool(
	agentQueryRepository repository.AgentQueryRepository,
	userId uuid.UUID,
	calendarId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*CalendarEventSearchTool, error) {
	if agentQueryRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &CalendarEventSearchTool{
		agentQueryRepository: agentQueryRepository,
		userId:               userId,
		calendarId:           calendarId,
		streamFn:             streamFn,
	}, nil
}

func (t CalendarEventSearchTool) Name() string {
	// return "予定検索ツール"
	return "search_calendar_events"
}

func (t CalendarEventSearchTool) Description() string {
	// return "期間を指定して予定を検索します。開始日時と終了日時をRFC3339形式で指定して検索します。"
	return "Searches for calendar events within a specified time range. Requires start and end times in RFC3339 format (e.g., 2025-01-01T09:00:00Z)."
}

func (t CalendarEventSearchTool) Strict() bool {
	return true
}

func (t CalendarEventSearchTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"start_at": map[string]any{
				"type":   "string",
				"format": "date-time",
			},
			"end_at": map[string]any{
				"type":   "string",
				"format": "date-time",
			},
		},
		"required": []string{"start_at", "end_at"},
	}
}

type CalendarEventSearchInput struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

func (t CalendarEventSearchTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Debug("calendar event search called")
	var inputModel EventSearchInput
	if err := json.Unmarshal([]byte(input), &inputModel); err != nil {
		return "", err
	}
	if inputModel.StartAt.IsZero() || inputModel.EndAt.IsZero() {
		logrus.WithField("user_input", input).Error("invalid user input time")
		return "", errors.New("input time is nil")
	}
	err := t.streamFn(ctx, statusChunk("fetching_events"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	events, err := t.agentQueryRepository.FindCalendarEventByUserIdAndDate(
		ctx,
		t.userId,
		t.calendarId,
		inputModel.StartAt,
		inputModel.EndAt,
	)
	if err != nil {
		return "", err
	}
	output, err := json.Marshal(events)
	if err != nil {
		return "", err
	}
	logrus.WithFields(logrus.Fields{
		"len":      len(events),
		"start_at": inputModel.StartAt.Format("2006-01-02 15:04:05"),
		"end_at":   inputModel.EndAt.Format("2006-01-02 15:04:05"),
	}).Debug("calendar event search tool called")
	return string(output), nil
}
