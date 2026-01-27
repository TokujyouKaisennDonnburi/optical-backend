package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SchedulerResultResponse struct {
	OwnerId   uuid.UUID        `json:"ownerId"`
	Title     string           `json:"title"`
	Memo      string           `json:"memo"`
	LimitTime time.Time        `json:"limitTime"`
	IsAllDay  bool             `json:"isAllDay"`
	Members   []MemberResponse `json:"members"`
	Date      []DateResponse   `json:"date"`
}
type MemberResponse struct {
	UserId   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
}
type DateResponse struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (h *SchedulerHttpHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	// schedulerId
	schedulerId, err := uuid.Parse(chi.URLParam(r, "schedulerId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	result, err := h.schedulerQuery.SchedulerResult(r.Context(), query.SchedulerResultInput{
		SchedulerId: schedulerId,
		UserId:      userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseMembers := make([]MemberResponse, len(result.Members))
	for i, v := range result.Members {
		responseMembers[i] = MemberResponse{
			UserId:   v.UserId,
			UserName: v.UserName,
		}
	}
	responseDates := make([]DateResponse, len(result.Date))
	for i, v := range result.Date {
		responseDates[i] = DateResponse{
			Date:      v.Date,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		}
	}
	responseResult := SchedulerResultResponse{
		OwnerId:   result.OwnerId,
		Title:     result.Title,
		Memo:      result.Memo,
		LimitTime: result.LimitTime,
		IsAllDay:  result.IsAllDay,
		Members:   responseMembers,
		Date:      responseDates,
	}
	// response
	render.JSON(w, r, responseResult)
}
