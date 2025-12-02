package auth

import (
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	USER_ID_CONTEXT_KEY = "userId"
)

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
func GetJwtSecretKey() []byte {
	secret, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		panic("\"JWT_SECRET_KEY\" is not set")
	}
	return []byte(secret)
}
