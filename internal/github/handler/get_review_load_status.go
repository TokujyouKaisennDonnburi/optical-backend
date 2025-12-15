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

type GithubReviewLoadStatusListResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Reviewed int    `json:"reviewed"`
	Assigned int    `json:"assigned"`
}

func (h *GithubHandler) GetReviewLoadStatus(w http.ResponseWriter, r *http.Request) {
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
	outputs, err := h.githubQuery.GetReviewLoadStatus(r.Context(), query.ReviewLoadStatusInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseList := make([]GithubReviewLoadStatusListResponse, len(outputs))
	for i, loadStatus := range outputs {
		responseList[i] = GithubReviewLoadStatusListResponse{
			Id:       loadStatus.GithubId,
			Name:     loadStatus.GithubName,
			Url:      loadStatus.GithubUrl,
			Reviewed: loadStatus.Reviewed,
			Assigned: loadStatus.Assigned,
		}
	}
	render.JSON(w, r, responseList)
}
