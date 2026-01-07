package tool

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OptionListTool struct {
	agentQueryRepository repository.AgentQueryRepository
	userId               uuid.UUID
	streamFn             func(context.Context, []byte) error
}

func NewOptionListTool(
	agentQueryRepository repository.AgentQueryRepository,
	userId uuid.UUID,
	streamFn func(context.Context, []byte) error,
) (*OptionListTool, error) {
	if agentQueryRepository == nil {
		return nil, errors.New("eventAgentRepository is nil")
	}
	return &OptionListTool{
		agentQueryRepository: agentQueryRepository,
		userId:               userId,
		streamFn:             streamFn,
	}, nil
}

func (t OptionListTool) Name() string {
	return "list_option"
}

func (t OptionListTool) Description() string {
	return "Retrieves all options can set to calendar."
}

func (t OptionListTool) Strict() bool {
	return false
}

type OptionListModel struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t OptionListTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t OptionListTool) Call(ctx context.Context, input string) (string, error) {
	logrus.WithField("user_input", input).Info("option list tool called")
	err := t.streamFn(ctx, statusChunk("fetching_options"))
	if err != nil {
		logrus.WithError(err).Error("progress streaming error")
	}
	options, err := t.agentQueryRepository.FindOptions(ctx)
	if err != nil {
		return "", err
	}
	output, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	logrus.WithField("len", len(options)).Info("event create tool called")
	return string(output), nil
}
