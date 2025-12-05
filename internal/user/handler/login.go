package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken  string                `json:"accessToken"`
	RefreshToken string                `json:"refreshToken"`
	User         UserLoginResponseUser `json:"user"`
}

type UserLoginResponseUser struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ユーザにログインする
func (h *UserHttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request UserLoginRequest
	// リクエストを取得
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		err = render.Render(w, r, apperr.ErrInvalidRequest(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	output, err := h.userCommand.LoginUser(r.Context(), command.UserLoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		err = render.Render(w, r, apperr.ErrUnauthorized(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	render.JSON(w, r, UserLoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User: UserLoginResponseUser{
			Id:    output.Id.String(),
			Name:  output.Name,
			Email: output.Email,
		},
	})
}
