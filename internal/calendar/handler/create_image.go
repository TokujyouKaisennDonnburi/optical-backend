package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SaveImageResponse struct {
	Id uuid.UUID `json:"id"`
}

func (h *CalendarHttpHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		err = render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	defer file.Close()
	output, err := h.calendarCommand.SaveImage(r.Context(), command.SaveImageCommandInput{
		File:   file,
		Header: header,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, SaveImageResponse{
		Id: output.Id,
	})
}
