package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreateResponse struct {
	Id           string `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ユーザーを新規作成する
func (h *UserHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request UserCreateRequest
	// リクエスト取得
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		err = render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	// ユーザー作成
	output, err := h.userCommand.CreateUser(context.Background(), command.UserCreateInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// レスポンス書き込み
	render.JSON(w, r, UserCreateResponse{
		Id:           output.Id.String(),
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}
