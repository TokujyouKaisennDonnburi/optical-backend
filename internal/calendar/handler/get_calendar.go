package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/google/uuid"
)

type CalendarResponse struct {
	Id      string                    `json:"id"`
	Name    string                    `json:"name"`
	Color   string                    `json:"color"`
	Image   string                    `json:"imageUrl"`
	Members []CalendarMemberResponse  `json:"member"`
	Options []CalendarOptionsResponse `json:"option"`
}
type CalendarMemberResponse struct {
	UserId   string    `json:"userId"`
	Name     string    `json:"name"`
	JoinedAt time.Time `json:"joinedAt"`
}
type CalendarOptionsResponse struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Deprecated bool   `json:"deprecated"`
}

func (h *CalendarHttpHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// calendarId
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// input
	output, err := h.calendarQuery.GetCalendar(r.Context(), query.GetCalendarInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	members := make([]CalendarMemberResponse, len(output.Members))
	for i, row := range output.Members {
		members[i] = CalendarMemberResponse{
			UserId:   row.UserId.String(),
			Name:     row.Name,
			JoinedAt: row.JoinedAt,
		}
	}
	options := make([]CalendarOptionsResponse, len(output.Options))
	for i, row := range output.Options {
		options[i] = CalendarOptionsResponse{
			Id:         row.Id,
			Name:       row.Name,
			Deprecated: row.Deprecated,
		}
	}
	response := CalendarResponse{
		Id:      output.Id.String(),
		Name:    output.Name,
		Color:   string(output.Color),
		Image:   output.ImageUrl,
		Members: members,
		Options: options,
	}
	render.JSON(w, r, response)
}
