package tool

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CalendarOptionUpdateTool struct {
	agentCommandRepository repository.AgentCommandRepository
	userId                 uuid.UUID
	calendarId             uuid.UUID
	streamFn               func(context.Context, []byte) error
}

func NewCalendarOptionUpdateTool(
	agentCommandRepository repository.AgentCommandRepository,
	userId, calendarId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*CalendarOptionUpdateTool, error) {
	if agentCommandRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &CalendarOptionUpdateTool{
		agentCommandRepository: agentCommandRepository,
		userId:                 userId,
		calendarId:             calendarId,
		streamFn:               streamFn,
	}, nil
}

func (t CalendarOptionUpdateTool) Name() string {
	return "update_calendar_option"
}

func (t CalendarOptionUpdateTool) Description() string {
	return "Overwrites the user's calendar  options. This replaces all specified settings entirely."
}

func (t CalendarOptionUpdateTool) Strict() bool {
	return false
}

type OptionUpdateInputModel struct {
	Options []int32 `json:"options"`
}

func (t CalendarOptionUpdateTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"options": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":   "integer",
					"format": "int32",
				},
			},
		},
		"required": []string{"options"},
	}
}

func (t CalendarOptionUpdateTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Info("calendar option update tool called")
	err := t.streamFn(ctx, statusChunk("updating_options"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	var model OptionUpdateInputModel
	if err := json.Unmarshal([]byte(input), &model); err != nil {
		return "false", err
	}
	err = t.agentCommandRepository.UpdateOptions(ctx, t.userId, t.calendarId, model.Options)
	if err != nil {
		return "false", err
	}
	logrus.WithField("len", len(model.Options)).Info("event create tool called")
	return "true", nil
}
