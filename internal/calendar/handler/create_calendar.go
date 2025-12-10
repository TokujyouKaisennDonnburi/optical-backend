package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CalendarCreateRequest struct {
	Name      string   `json:"name"`
	Color     string   `json:"color"`
	ImageId   string   `json:"imageId"`
	Members   []string `json:"members"`
	OptionIds []int32  `json:"optionIds"`
}

type CalendarCreateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *CalendarHttpHandler) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	userName, err := auth.GetUserNameFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	var request CalendarCreateRequest
	// リクエストJSONをバインド
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	var imageId uuid.UUID
	if request.ImageId != "" {
		imageId, err = uuid.Parse(request.ImageId)
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
			return
		}
	}
	// カレンダーを作成
	output, err := h.calendarCommand.CreateCalendar(context.Background(), command.CalendarCreateInput{
		UserId:        userId,
		UserName:      userName,
		ImageId:       imageId,
		CalendarName:  request.Name,
		CalendarColor: request.Color,
		MemberEmails:  request.Members,
		OptionIds:     request.OptionIds,
	})
	if err != nil {
		apperr.HandleAppError(w,r,err)
		return
	}
	// レスポンスに変換
	render.JSON(w, r, CalendarCreateResponse{
		Id:   output.Id.String(),
		Name: output.Name,
	})
}
