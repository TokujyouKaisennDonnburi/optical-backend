package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type UserAuthMiddleware struct{}

func NewUserAuthMiddleware() *UserAuthMiddleware {
	return &UserAuthMiddleware{}
}

// JWT認証を行うミドルウェア
func (m *UserAuthMiddleware) JWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ヘッダーからトークンを取得
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			_ = render.Render(w, r, apperr.ErrUnauthorized(errors.New("Authorization header is invalid")))
			return
		}
		authorizationHeader = strings.TrimPrefix(authorizationHeader, "Bearer ")
		// トークンデコード
		accesToken, err := user.DecodeAccessToken(authorizationHeader)
		if err != nil {
			_ = render.Render(w, r, apperr.ErrUnauthorized(err))
			return
		}
		if accesToken.IsExpired() {
			_ = render.Render(w, r, apperr.ErrUnauthorized(errors.New("accessToken is expired")))
			return
		}
		// コンテキストに含めてエンドポイントに渡す
		ctx := context.WithValue(r.Context(), auth.USER_ID_CONTEXT_KEY, accesToken.UserId.String())
		ctx = context.WithValue(ctx, auth.USER_NAME_CONTEXT_KEY, accesToken.UserName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
