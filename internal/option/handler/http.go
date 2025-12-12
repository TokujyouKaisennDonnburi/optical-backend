package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/option/service/query"

type OptionHttpHandler struct {
	optionQuery *query.OptionQuery
}

func NewOptionHttpHandler(optionQuery *query.OptionQuery) *OptionHttpHandler {
	if optionQuery == nil {
		panic("OptionQuery is nil")
	}
	return &OptionHttpHandler{
		optionQuery: optionQuery,
	}
}

