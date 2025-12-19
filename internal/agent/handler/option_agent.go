package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type SuggestOptionRequest struct {
	Message string `json:"message"`
}

type SuggestOptionResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func (h *AgentHandler) SuggestOptions(w http.ResponseWriter, r *http.Request) {
	var request SuggestOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	outputs, err := h.agentQuery.SuggestOptions(r.Context(), query.SuggestOptionsInput{
		Message: request.Message,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseList := make([]SuggestOptionResponse, len(outputs))
	for i, output := range outputs {
		responseList[i] = SuggestOptionResponse{
			Id:   output.Id,
			Name: output.Name,
		}
	}
	render.JSON(w, r, responseList)
}
