package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type UserResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarUrl *string   `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 自身のユーザー情報を取得する
func (h *UserHttpHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// ユーザーIDを取得
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// ユーザー情報を取得
	output, err := h.userQuery.GetUser(r.Context(), query.UserQueryInput{
		UserId: userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	response := UserResponse{
		Id:        output.Id.String(),
		Name:      output.Name,
		Email:     output.Email,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	}
	if output.Avatar.Valid {
		response.AvatarUrl = &output.Avatar.Url
	}
	render.JSON(w, r, response)
}
