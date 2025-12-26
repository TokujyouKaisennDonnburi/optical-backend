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

type EventSearchTool struct {
	agentQueryRepository repository.AgentQueryRepository
	userId               uuid.UUID
	streamFn             func(context.Context, []byte) error
}

func NewEventSearchTool(
	agentQueryRepository repository.AgentQueryRepository,
	userId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*EventSearchTool, error) {
	if agentQueryRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &EventSearchTool{
		agentQueryRepository: agentQueryRepository,
		userId:               userId,
		streamFn:             streamFn,
	}, nil
}

func (t EventSearchTool) Name() string {
	// return "予定検索ツール"
	return "search_events"
}

func (t EventSearchTool) Description() string {
	// return "期間を指定して予定を検索します。開始日時と終了日時をRFC3339形式で指定して検索します。"
	return "Searches for calendar events within a specified time range. Requires start and end times in RFC3339 format (e.g., 2025-01-01T09:00:00Z)."
}

func (t EventSearchTool) Strict() bool {
	return true
}

func (t EventSearchTool) Parameters() map[string]any {
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

type EventSearchInput struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

func (t EventSearchTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Debug("event search called")
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
	events, err := t.agentQueryRepository.FindEventByUserIdAndDate(ctx, t.userId, inputModel.StartAt, inputModel.EndAt)
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
	}).Debug("event search tool called")
	return string(output), nil
}
