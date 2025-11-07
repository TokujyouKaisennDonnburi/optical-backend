package handler

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	USER_ID_CONTEXT_KEY = "userId"
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
			return getJwtSecretKey(), nil
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
		ctx := context.WithValue(r.Context(), USER_ID_CONTEXT_KEY, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIdFromContext(r *http.Request) (uuid.UUID, error) {
	// IDを取得
 	uid := r.Context().Value(USER_ID_CONTEXT_KEY)
	uidStr, ok := uid.(string)
	if !ok {
		return uuid.Nil, errors.New("invalid jwt context")
	}
	return uuid.Parse(uidStr)
}

// JWTの暗号化鍵
func getJwtSecretKey() []byte {
	secret, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		panic("\"JWT_SECRET_KEY\" is not set")
	}
	return []byte(secret)
}
