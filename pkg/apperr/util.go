package apperr

import "net/http"

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	apperr, ok := err.(*AppError)
	if !ok {
		return false
	}
	return apperr.StatusCode == http.StatusNotFound
}

func IsForbidden(err error) bool {
	if err == nil {
		return false
	}
	apperr, ok := err.(*AppError)
	if !ok {
		return false
	}
	return apperr.StatusCode == http.StatusForbidden
}
