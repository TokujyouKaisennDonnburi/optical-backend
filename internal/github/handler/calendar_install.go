package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type GithubCalendarInstallRequest struct {
	CalendarId     string `json:"calendarId"`
	InstallationId string `json:"installationId"`
}

func (h *GithubHandler) InstallToCalendar(w http.ResponseWriter, r *http.Request) {
	var request GithubCalendarInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	err := h.githubCommand.InstallToCalendar(r.Context(), command.GithubCalendarInstallInput{
		CalendarId:     request.CalendarId,
		InstallationId: request.InstallationId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.NoContent(w, r)
}
