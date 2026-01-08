package tool

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type EventCreateTool struct {
	agentCommandRepository repository.AgentCommandRepository
	userId                 uuid.UUID
	streamFn               func(context.Context, []byte) error
}

func NewEventCreateTool(
	agentCommandRepository repository.AgentCommandRepository,
	userId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*EventCreateTool, error) {
	if agentCommandRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &EventCreateTool{
		agentCommandRepository: agentCommandRepository,
		userId:                 userId,
		streamFn:               streamFn,
	}, nil
}

func (t EventCreateTool) Name() string {
	// "予定一括作成ツール"
	return "bulk_create_events"
}

func (t EventCreateTool) Description() string {
	// カレンダーに予定を一括作成します。作成された予定を出力します。
	return "Creates multiple events within a single specified calendar in one operation. returns the details of the created events."
}

func (t EventCreateTool) Strict() bool {
	return false
}

type UserInputModel struct {
	CalendarId uuid.UUID          `json:"calendar_id"`
	Events     []EventCreateModel `json:"events"`
}

type EventCreateModel struct {
	Title      string    `json:"event_title"`
	Memo       string    `json:"memo"`
	Location   string    `json:"location"`
	IsAllday   bool      `json:"is_allday"`
	StartAt    time.Time `json:"start_at"`
	EndAt      time.Time `json:"end_at"`
}

func (t EventCreateTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"calendar_id": map[string]any{
				"type":   "string",
				"format": "uuid",
			},
			"events": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"event_title": map[string]any{
							"type": "string",
						},
						"memo": map[string]any{
							"type": "string",
						},
						"location": map[string]any{
							"type": "string",
						},
						"is_allday": map[string]any{
							"type": "boolean",
						},
						"start_at": map[string]any{
							"type":        "string",
							"description": "RFC3339",
							"format":      "date-time",
						},
						"end_at": map[string]any{
							"type":        "string",
							"description": "RFC3339",
							"format":      "date-time",
						},
					},
					"required": []string{"event_title", "is_allday", "start_at", "end_at"},
					"additionalProperties": false,
				},
			},
		},
	}
}

func (t EventCreateTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Info("event create tool called")
	err := t.streamFn(ctx, statusChunk("creating_events"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	var model UserInputModel
	if err := json.Unmarshal([]byte(input), &model); err != nil {
		return "", err
	}
	if len(model.Events) == 0 {
		return "", nil
	}
	events := make([]calendar.Event, len(model.Events))
	analyzableModels := make([]agent.AnalyzableEvent, len(model.Events))
	err = t.agentCommandRepository.CreateEvents(
		ctx,
		t.userId,
		model.CalendarId,
		func(c *calendar.Calendar) ([]calendar.Event, error) {
			isMember := false
			for _, member := range c.Members {
				if member.UserId == t.userId {
					isMember = true
					break
				}
			}
			if !isMember {
				return nil, errors.New("invalid user access")
			}
			for i, event := range model.Events {
				scheduledTime, err := calendar.NewScheduledTime(event.IsAllday, event.StartAt, event.EndAt)
				if err != nil {
					return nil, err
				}
				newEvent, err := calendar.NewEvent(c.Id, event.Title, event.Memo, string(c.Color), event.Location, *scheduledTime)
				if err != nil {
					return nil, err
				}
				events[i] = *newEvent
				analyzableModels[i] = agent.AnalyzableEvent{
					CalendarId:    c.Id.String(),
					CalendarName:  c.Name,
					CalendarColor: string(c.Color),
					Id:            newEvent.Id.String(),
					Title:         newEvent.Title,
					Memo:          newEvent.Memo,
					Location:      newEvent.Location,
					StartAt:       event.StartAt,
					EndAt:         event.EndAt,
					IsAllday:      event.IsAllday,
				}
			}
			return events, nil
		},
	)
	if err != nil {
		return "", err
	}
	output, err := json.Marshal(analyzableModels)
	if err != nil {
		return "", err
	}
	logrus.WithField("len", len(events)).Info("event create tool called")
	return string(output), nil
}
