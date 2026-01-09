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
	return "list_events"
}

func (t CalendarEventSearchTool) Description() string {
	// return "Searches for calendar events within a specified time range. Requires start and end times in "
	return "Searches for calendar events based on specified criteria such as title, location, and time range. Use this to find existing appointments or check availability for specific dates."
}

func (t CalendarEventSearchTool) Strict() bool {
	return false
}

func (t CalendarEventSearchTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "The keywords to search for in event titles",
			},
			"location": map[string]any{
				"type":        "string",
				"description": "event location",
			},
			"start_at": map[string]any{
				"type":        "string",
				"format":      "date-time",
				"description": "RFC3339 format (e.g., 2025-01-01T09:00:00Z)",
			},
			"end_at": map[string]any{
				"type":        "string",
				"format":      "date-time",
				"description": "RFC3339 format (e.g., 2025-01-01T09:00:00Z)",
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "Maximum number of results to return.",
				"minimum":     0,
				"maximum":     65535,
			},
		},
		"required": []string{"start_at", "end_at"},
	}
}

type CalendarEventSearchInput struct {
	Title    string    `json:"title"`
	Location string    `json:"location"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	Limit    uint16    `json:"limit"`
}

func (t CalendarEventSearchTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Info("calendar event search called")
	var inputModel CalendarEventSearchInput
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
		inputModel.Title,
		inputModel.Location,
		inputModel.StartAt,
		inputModel.EndAt,
		inputModel.Limit,
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
		"title":    inputModel.Title,
		"location": inputModel.Location,
		"start_at": inputModel.StartAt.Format("2006-01-02 15:04:05"),
		"end_at":   inputModel.EndAt.Format("2006-01-02 15:04:05"),
		"limit":    inputModel.Limit,
	}).Info("calendar event search tool called")
	return string(output), nil
}
