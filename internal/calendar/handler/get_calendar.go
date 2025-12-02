package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
)


type CalendarResponse struct {
	Id		string `json:"id"`
	Name    string
	Color   string
	Image   string
	Members []calrendar.
	Options []option.Option
}
