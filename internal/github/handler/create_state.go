package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type GithubStateResponse struct {
	Url string `json:"url"`
}

func (h *GithubHandler) CreateState(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	output, err := h.githubCommand.CreateOauthState(r.Context(), command.GithubOauthStateInput{
		UserId: userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, GithubStateResponse{
		Url: output.Url,
	})
}
