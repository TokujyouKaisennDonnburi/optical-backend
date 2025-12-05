package handler

import (
	"errors"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	ErrInvalidCode = errors.New("invalid code")
	ErrInvalidInstallationId = errors.New("invalid installation id")
)

func (h *UserHttpHandler) LinkGithub(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		err := render.Render(w, r, apperr.ErrInvalidRequest(ErrInvalidCode))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	installationId := chi.URLParam(r, "installation_id")
	if installationId == "" {
		err := render.Render(w, r, apperr.ErrInvalidRequest(ErrInvalidInstallationId))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
}
