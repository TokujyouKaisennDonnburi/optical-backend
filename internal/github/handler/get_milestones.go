package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type GithubMilestonesRequest struct {
	UserId     uuid.UUID `json:"userId"`
	CalendarId uuid.UUID `json:"calendarId"`
}

type GithubMilestonesResponse struct {
	Title    string `json:"title"`
	Progress int    `json:"progress"`
	Open     int    `json:"open"`
	Close    int    `json:"close"`
}

func (h *GithubHandler) GetMilestones(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	milestones, err := h.githubQuery.GetMilestone(r.Context(), query.MilestonesInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseList := make([]GithubMilestonesResponse, len(milestones))
	for i, milestone := range milestones {
		responseList[i] = GithubMilestonesResponse{
			Title:    milestone.Title,
			Progress: milestone.Progress,
			Open:     milestone.Open,
			Close:    milestone.Close,
		}
	}
	render.JSON(w, r, responseList)
}
