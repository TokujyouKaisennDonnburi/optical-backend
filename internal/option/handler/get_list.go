package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type OptionResponse struct {
	OptionId   string `json:"optionId"`
	Name       string `json:"name"`
	Deprecated bool   `json:"deprecated"`
}

func (h *OptionHttpHandler) GetList(w http.ResponseWriter, r *http.Request) {
	// userId
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	output, err := h.optionQuery.GetListOption(r.Context(), query.ListQueryInput{
		UserId: userId,
	})
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, OptionResponse{
		Id:         output.Id,
		Name:       output.Name,
		Deprecated: output.Deprecated,
	})
}

