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
	RepositoryId   int64                                    `json:"repositoryId"`
	RepositoryName string                                   `json:"repositoryName"`
	Total          int                                      `json:"total"`
	Reviewers      []GithubReviewLoadStatusListResponseItem `json:"reviewers"`
}

type GithubReviewLoadStatusListResponseItem struct {
	GithubId   int64  `json:"githubId"`
	GithubName string `json:"githubName"`
	GithubUrl  string `json:"githubUrl"`
	Assigned   int    `json:"assigned"`
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
		reviewers := make([]GithubReviewLoadStatusListResponseItem, len(loadStatus.Status))
		for i, status := range loadStatus.Status {
			reviewers[i] = GithubReviewLoadStatusListResponseItem{
				GithubId:   status.GithubId,
				GithubName: status.GithubName,
				GithubUrl:  status.GithubUrl,
				Assigned:   status.Assigned,
			}
		}
		responseList[i] = GithubReviewLoadStatusListResponse{
			RepositoryId:   loadStatus.RepositoryId,
			RepositoryName: loadStatus.RepositoryName,
			Total:          loadStatus.Total,
			Reviewers:      reviewers,
		}
	}
	render.JSON(w, r, responseList)
}
