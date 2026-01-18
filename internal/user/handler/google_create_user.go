package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CreateGoogleUserRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type CreateGoogleUserResponse struct {
	UserId       uuid.UUID `json:"user"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

// Googleアカウントからユーザーを新規作成
func (h *UserHttpHandler) CreateGoogleUser(w http.ResponseWriter, r *http.Request) {
	var request CreateGoogleUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	output, err := h.userCommand.CreateGoogleUser(r.Context(), command.CreateGoogleUserInput{
		Code:  request.Code,
		State: request.State,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, CreateGoogleUserResponse{
		UserId:       output.UserId,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}
