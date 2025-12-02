package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type UserAuthMiddleware struct {}

func NewUserAuthMiddleware() *UserAuthMiddleware {
	return &UserAuthMiddleware{}
}

// JWT認証を行うミドルウェア
func (m *UserAuthMiddleware) JWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ヘッダーからトークンを取得
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			err := render.Render(w, r, apperr.ErrUnauthorized(errors.New("Authorization header is invalid")))
			if err != nil {
				_ = render.Render(w, r, apperr.ErrInternalServerError(err))
			}
			return
		}
		authorizationHeader = strings.TrimPrefix(authorizationHeader, "Bearer ")
		// トークンデコード
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authorizationHeader, claims, func(t *jwt.Token) (any, error) {
			return auth.GetJwtSecretKey(), nil
		})
		if err != nil {
			err = render.Render(w, r, apperr.ErrUnauthorized(err))
			if err != nil {
				_ = render.Render(w, r, apperr.ErrInternalServerError(err))
			}
			return
		}
		// ユーザーIDを取得
		userId, err := claims.GetSubject()
		if err != nil {
			err = render.Render(w, r, apperr.ErrUnauthorized(err))
			if err != nil {
				_ = render.Render(w, r, apperr.ErrInternalServerError(err))
			}
		}
		// コンテキストに含めてエンドポイントに渡す
		ctx := context.WithValue(r.Context(), auth.USER_ID_CONTEXT_KEY, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

