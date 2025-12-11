package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type GithubUserLinkRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func (h *GithubHandler) LinkUser(w http.ResponseWriter, r *http.Request) {
	var request GithubUserLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	err := h.githubCommand.LinkUser(r.Context(), command.GithubLinkUserInput{
		Code:  request.Code,
		State: request.State,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.NoContent(w, r)
}
