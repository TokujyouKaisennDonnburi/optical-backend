package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// 未連携時用でomitemptyを使用
type IsInstalledGithubAppResponse struct {
	IsInstalled    bool   `json:"isInstalled"`
	GithubId       string `json:"githubId,omitempty"`
	GithubName     string `json:"githubName,omitempty"`
	InstallationId string `json:"installationId,omitempty"`
	InstalledAt    string `json:"installedAt,omitempty"`
}

func (h *GithubHandler) IsInstalledGithubApp(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// calendarId取得
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// handlerでquery呼び出し
	result, err := h.githubQuery.IsInstalledGithubApp(r.Context(), userId, calendarId)
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	var installedAt string
	if !result.InstalledAt.IsZero() {
		installedAt = result.InstalledAt.Format(time.RFC3339)
	}

	render.JSON(w, r, IsInstalledGithubAppResponse{
		IsInstalled:    result.IsInstalled,
		GithubId:       result.GithubId,
		GithubName:     result.GithubName,
		InstallationId: result.InstallationId,
		InstalledAt:    installedAt,
	})
}
