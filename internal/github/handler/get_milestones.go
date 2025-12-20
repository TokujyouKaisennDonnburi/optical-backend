package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type GithubMilestonesRequest struct {
	UserId    uuid.UUID `json:"userId"`
	ClendarId uuid.UUID `json:"calendarId"`
}

type GithubMilestonesResponse struct {
	Title    string `json:"title"`
	Progress int8   `json:"progress"`
	Open     int8   `json:"open"`
	Close    int8   `json:"close"`
}

func (h *GithubHandler) GetMilestones(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetClientId("userId")
	if err != nil {
		return _  = render.Render(w, r, apperr.ErrInternalServerError(err))
	}
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		return _ = render.Render(w, r, apperr.ErrInvalidRequest(err))
	}

}
