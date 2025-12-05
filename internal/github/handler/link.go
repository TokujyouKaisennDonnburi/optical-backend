package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	// "github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var (
	ErrInvalidCode           = errors.New("invalid code")
	ErrInvalidInstallationId = errors.New("invalid installation id")
)

type GithubLinkRequest struct {
	Code           string `json:"code"`
	InstallationId string `json:"installationId"`
}

func (h *GithubHandler) LinkGithub(w http.ResponseWriter, r *http.Request) {
	userId := uuid.New()
	// userId, err := handler.GetUserIdFromContext(r)
	// if err != nil {
	// 	err := render.Render(w, r, apperr.ErrInvalidRequest(err))
	// 	if err != nil {
	// 		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
	// 	}
	// 	return
	// }
	var request GithubLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		err := render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
	}
	err := h.githubCommand.LinkUser(r.Context(), command.GithubLinkInput{
		UserId:         userId,
		InstallationId: request.InstallationId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.NoContent(w, r)
}
