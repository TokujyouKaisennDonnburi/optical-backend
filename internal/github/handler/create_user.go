package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type GithubCreateUserOauthStateResponse struct {
	Url string `json:"url"`
}

func (h *GithubHandler) CreateNewUserOauthState(w http.ResponseWriter, r *http.Request) {
	output, err := h.githubCommand.CreateOauthState(r.Context(), command.GithubOauthStateInput{
		UserId: uuid.Max,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, GithubOauthStateResponse{
		Url: output.Url,
	})
}
