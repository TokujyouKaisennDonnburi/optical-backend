package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CreateGoogleOauthStateResponse struct {
	Url string `json:"url"`
}

func (h *UserHttpHandler) CreateGoogleOauthState(w http.ResponseWriter, r *http.Request) {
	output, err := h.userCommand.CreateGoogleState(r.Context(), command.CreateGoogleStateInput{
		UserId: uuid.Max,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, CreateGoogleOauthStateResponse{
		Url: output.Url,
	})
}
