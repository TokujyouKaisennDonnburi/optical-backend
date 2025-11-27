package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CalendarCreateRequest struct {
	Name      string   `json:"name"`
	Color     string   `json:"color"`
	ImageId   string   `json:"imageId"`
	OptionIds []string `json:"optionIds"`
}

type CalendarCreateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *CalendarHttpHandler) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	var request CalendarCreateRequest
	// リクエストJSONをバインド
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		err = render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	optionIds := []uuid.UUID{}
	for _, id := range request.OptionIds {
		optionId, err := uuid.Parse(id)
		if err != nil {
			err = render.Render(w, r, apperr.ErrInvalidRequest(err))
			if err != nil {
				_ = render.Render(w, r, apperr.ErrInternalServerError(err))
			}
			return
		}
		optionIds = append(optionIds, optionId)
	}
	imageId, err := uuid.Parse(request.ImageId)
	if err != nil {
		err = render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	// カレンダーを作成
	output, err := h.calendarCommand.CreateCalendar(context.Background(), command.CalendarCreateArgs{
		UserId:        userId,
		ImageId:       imageId,
		CalendarName:  request.Name,
		CalendarColor: request.Color,
		OptionIds:     optionIds,
	})
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// レスポンスに変換
	render.JSON(w, r, CalendarCreateResponse{
		Id:   output.Id.String(),
		Name: output.Name,
	})
}
