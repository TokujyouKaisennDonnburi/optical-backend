package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
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
	calendarId, err := handler.GetCalendarIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	
	response := CalendarResponse{
		Id:
		Name:
		Color:
		Image:
		Members:
		Options:
	}
	render.JSON(w,r,response)
}
