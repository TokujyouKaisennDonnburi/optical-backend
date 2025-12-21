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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AgentCalendarChatRequest struct {
	Message string `json:"message"`
}

func (h *AgentHandler) CalendarChat(w http.ResponseWriter, r *http.Request) {
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
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		logrus.WithError(err).Error("userId decrypt error")
		return
	}

	var request AgentCalendarChatRequest
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

	err = h.agentCommand.CalendarChat(r.Context(), command.AgentCommandCalendarChatInput{
		UserId:      userId,
		CalendarId:  calendarId,
		UserInput:   request.Message,
		StreamingFn: streamingFn,
	})
	if err != nil {
		logrus.WithError(err).Error("agent exec error")
		return
	}
}

