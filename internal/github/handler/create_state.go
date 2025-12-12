package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type GithubAppStateRequest struct {
	CalendarId uuid.UUID
}

type GithubAppStateResponse struct {
	Url string `json:"url"`
}

func (h *GithubHandler) CreateAppState(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	var request GithubAppStateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	output, err := h.githubCommand.CreateAppState(r.Context(), command.GithubAppStateInput{
		UserId:     userId,
		CalendarId: request.CalendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, GithubAppStateResponse{
		Url: output.Url,
	})
}

type GithubOauthStateResponse struct {
	Url string `json:"url"`
}

func (h *GithubHandler) CreateOauthState(w http.ResponseWriter, r *http.Request) {
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
	render.JSON(w, r, GithubOauthStateResponse{
		Url: output.Url,
	})
}
