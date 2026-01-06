package tool

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CalendarOptionListTool struct {
	agentQueryRepository repository.AgentQueryRepository
	userId               uuid.UUID
	calendarId           uuid.UUID
	streamFn             func(context.Context, []byte) error
}

func NewCalendarOptionListTool(
	agentQueryRepository repository.AgentQueryRepository,
	userId, calendarId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*CalendarOptionListTool, error) {
	if agentQueryRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &CalendarOptionListTool{
		agentQueryRepository: agentQueryRepository,
		userId:               userId,
		calendarId:           calendarId,
		streamFn:             streamFn,
	}, nil
}

func (t CalendarOptionListTool) Name() string {
	return "list_calendar_option"
}

func (t CalendarOptionListTool) Description() string {
	return "Retrieves all current options configured in the user's calendar."
}

func getDescriptionMap() map[int32]string {
	return map[int32]string{
		1: "Githubでユーザーがレビューを依頼されているものを確認できます",
		2: "Githubのどのユーザーにレビュー依頼がされているかの付加状況を確認できます",
	}
}

func (t CalendarOptionListTool) Strict() bool {
	return false
}

type CalendarOptionListModel struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t CalendarOptionListTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t CalendarOptionListTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Info("calendar option list tool called")
	err := t.streamFn(ctx, statusChunk("fetching_options"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	options, err := t.agentQueryRepository.FindOptionsByCalendarId(ctx, t.userId, t.calendarId)
	if err != nil {
		return "", err
	}
	outputList := make([]CalendarOptionListModel, len(options))
	for i, option := range options {
		outputList[i] = CalendarOptionListModel{
			Id:          option.Id,
			Name:        option.Name,
			Description: getDescriptionMap()[option.Id],
		}
	}
	output, err := json.Marshal(outputList)
	if err != nil {
		return "", err
	}
	logrus.WithField("len", len(outputList)).Info("event create tool called")
	return string(output), nil
}
