package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenRefreshResponse struct {
	AccessToken string    `json:"accessToken"`
	ExpiresIn   time.Time `json:"expiresIn"`
}

// トークンリフレッシュ用エンドポイント
func (h UserHttpHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request TokenRefreshRequest
	// リクエスト取得
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		err = render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	// トークンをリフレッシュ
	output, err := h.userCommand.Refresh(r.Context(), command.TokenRefreshInput{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		err = render.Render(w, r, apperr.ErrUnauthorized(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	// レスポンスを返す
	render.JSON(w, r, TokenRefreshResponse{
		AccessToken: output.AccessToken,
		ExpiresIn:   output.ExpiresIn,
	})
}
