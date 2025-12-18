package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type UploadAvatarResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (h *UserHttpHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	defer file.Close()
	output, err := h.userCommand.UploadAvatar(r.Context(), command.UploadAvatarInput{
		UserId: userId,
		File:   file,
		Header: header,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, UploadAvatarResponse{
		Id:  output.Id.String(),
		Url: output.Url,
	})
}
