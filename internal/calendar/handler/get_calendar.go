package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type CalendarResponse struct {
	Id		string `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Image   string `json:"image"`
	Members []calendar.Member `json:"member"`
	Options []option.Option   `json:"option"`
}

func (h *CalendarHttpHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// get calendarId
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// input
	output, err := h.calendarQuery.GetCalendar(r.Context(), query.CalendarQueryInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// output
	response := CalendarResponse{
		Id:		 string(output.Id),
		Name:	 output.Name,
		Color:	 output.Color,
		Image:	 output.Image.Url,
		Members: output.Members,
		Options: output.Options,
	}
	render.JSON(w,r,response)
}
