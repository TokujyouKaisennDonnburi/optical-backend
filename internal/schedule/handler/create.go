package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ScheduleCreateRequest struct {
	Name      string   `json:"name"`
	OptionIds []string `json:"optionIds"`
}

type ScheduleCreateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *ScheduleHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.Context().Value("id").(string))
	if err != nil {
		render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	var request ScheduleCreateRequest
	// リクエストJSONをバインド
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	optionIds := []uuid.UUID{}
	for _, id := range request.OptionIds {
		optionId, err := uuid.Parse(id)
		if err != nil {
			render.Render(w, r, apperr.ErrInvalidRequest(err))
			return
		}
		optionIds = append(optionIds, optionId)
	}
	// スケジュールを作成
	output, err := h.createCommandService.CreateSchedule(context.Background(), command.ScheduleCreateArgs{
		UserId:       userId,
		ScheduleName: request.Name,
		OptionIds:    optionIds,
	})
	if err != nil {
		render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// レスポンスに変換
	body, err := json.Marshal(ScheduleCreateResponse{
		Id:   output.Id.String(),
		Name: output.Name,
	})
	if err != nil {
		render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	w.Write(body)
}
