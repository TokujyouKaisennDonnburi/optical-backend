package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
)

type CalendarResponse struct {
	Id		string `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Image   string `json:"image"`
	Members []calendar.Member `json:"member"`
	Options []option.Option     `json:"option"`
}
