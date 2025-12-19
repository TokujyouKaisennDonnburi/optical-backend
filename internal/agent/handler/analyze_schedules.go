package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type AnalyzeSchedulesRequest struct {
	Request string `json:"request"`
}

func (h *AgentHandler) AnalyzeSchedules(w http.ResponseWriter, r *http.Request) {
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

	var request ExecAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.WithError(err).Error("agent message unmarshal error")
		return
	}
	streamingFn := func(ctx context.Context, b []byte) error {
		fmt.Fprintf(w, "data: {\"content\":\"%s\"}\n\n", string(b))
		flusher.Flush()
		return nil
	}

	err = h.agentQuery.AnalyzeSchedules(r.Context(), query.AnalyzeSchedulesQueryInput{
		UserId:      userId,
		UserInput:   request.Request,
		StreamingFn: streamingFn,
	})
	if err != nil {
		logrus.WithError(err).Error("agent exec error")
		return
	}
}
