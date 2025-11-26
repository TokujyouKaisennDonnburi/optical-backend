package apperr

import (
	"net/http"

	"github.com/go-chi/render"
)

type AppError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func (er *AppError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, er.StatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message: err.Error(),
	}
}

func ErrUnauthorized(err error) render.Renderer {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message: err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message: err.Error(),
	}
}
