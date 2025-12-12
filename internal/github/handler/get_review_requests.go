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

type GithubReviewRequestListResponse struct {
	Id        int64                                 `json:"id"`
	Url       string                                `json:"url"`
	Title     string                                `json:"title"`
	Number    int                                   `json:"number"`
	Assigness []GithubReviewRequestListResponseUser `json:"assigness"`
}

type GithubReviewRequestListResponseUser struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	AvatarUrl string `json:"avatarUrl"`
}

func (h *GithubHandler) GetReviewRequests(w http.ResponseWriter, r *http.Request) {
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
	pullRequests, err := h.githubQuery.GetReviewRequests(r.Context(), query.GithubReviewRequestsQueryInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseList := make([]GithubReviewRequestListResponse, len(pullRequests))
	for i, pullRequest := range pullRequests {
		assigness := make([]GithubReviewRequestListResponseUser, len(pullRequest.Assigness))
		for i, user := range pullRequest.Assigness {
			assigness[i] = GithubReviewRequestListResponseUser{
				Id:        user.Id,
				Name:      user.Login,
				Url:       user.Url,
				AvatarUrl: user.AvatarUrl,
			}
		}
		responseList[i] = GithubReviewRequestListResponse{
			Id:        pullRequest.Id,
			Url:       pullRequest.Url,
			Title:     pullRequest.Title,
			Number:    pullRequest.Number,
			Assigness: assigness,
		}
	}
	render.JSON(w, r, responseList)
}
