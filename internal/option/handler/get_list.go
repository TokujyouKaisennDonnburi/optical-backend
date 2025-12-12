package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type OptionResponse struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Deprecated bool   `json:"deprecated"`
}

func (h *OptionHttpHandler) GetList(w http.ResponseWriter, r *http.Request) {
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	output, err := h.optionQuery.GetListOption(r.Context(), query.OptionListQueryInput{
		UserId: userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	options := make([]OptionResponse, len(output))
	for i, opt := range output {
		options[i] = OptionResponse{
			Id:         opt.Id,
			Name:       opt.Name,
			Deprecated: opt.Deprecated,
		}
	}
	render.JSON(w, r, options)
}
