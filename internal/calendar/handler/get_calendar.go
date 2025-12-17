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
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Image   string   `json:"imageUrl"`
	Members []string `json:"member"`
	Options []string `json:"option"`
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
	// for文でループ
	response := make([]CalendarResponse, 0, len(output))
	for _, hoge := range output {
		response = append(response, CalendarResponse{
			Id:    output.Id.String(),
			Name:  output.Name,
			Color: string(output.Color),
			Image: output.Image.Url,
			Members: CalendarMemberResponse{
				UserId:   output.Id.String(),
				Name:     output.Name,
				JoinedAt: output.Members,
			},
			Options: CalendarOptionsResponse{
				Id:         append(output.Options),
				Name:       output.Options,
				Deprecated: output.Options,
			},
		})
	}
	render.JSON(w, r, response)
}
