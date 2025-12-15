package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type UserUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHttpHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	var request UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	err = h.userCommand.UpdateUser(r.Context(), command.UserUpdateInput{
		UserId: userId,
		Name:   request.Name,
		Email:  request.Email,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
