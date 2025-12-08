package apperr

import (
	"net/http"

	"github.com/go-chi/render"
)

func HandleAppError(w http.ResponseWriter, r *http.Request, err error) {
	apperr, ok := err.(*AppError)
	if !ok {
		_ = render.Render(w, r, ErrInternalServerError(err))
		return
	}
	err = render.Render(w, r, apperr)
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func InternalServerError(message string) error {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func ValidationError(message string) error {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func UnauthorizedError(message string) error {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func ForbiddenError(message string) error {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Message:    message,
	}
}

func NotFoundError(message string) error {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}
