package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type AgentChatRequest struct {
	Message string `json:"message"`
}

func (h *AgentHandler) Chat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(errors.New("streaming is not supported")))
		logrus.Error("flusher is not supported")
		return
	}
	defer func() {
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}()

	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		logrus.WithError(err).Error("userId decrypt error")
		return
	}

	var request AgentChatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.WithError(err).Error("agent message unmarshal error")
		return
	}
	streamingFn := func(ctx context.Context, b []byte) error {
		chunk := strings.ReplaceAll(string(b), "\n", "\\n")
		fmt.Fprintf(w, "data: %s\n\n", chunk)
		flusher.Flush()
		return nil
	}

	err = h.agentCommand.Chat(r.Context(), command.AgentCommandChatInput{
		UserId:      userId,
		UserInput:   request.Message,
		StreamingFn: streamingFn,
	})
	if err != nil {
		logrus.WithError(err).Error("agent exec error")
		jsonMessage, err := json.Marshal(err.Error())
		if err != nil {
			logrus.WithError(err).Error("agent exec json marshal error")
			return
		}
		errorMessage := `{"error":%s}`
		fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf(errorMessage, jsonMessage))
		return
	}
}
