package user

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ACCESS_TOKEN_EXPIRE  = 60 * 60 * 1        // 1時間
	REFRESH_TOKEN_EXPIRE = 60 * 60 * 24 * 180 // 180時間
)

type AccessToken struct {
	UserId    uuid.UUID
	Token     string
	ExpiresIn time.Time
}

func NewAccessToken(user *User) (*AccessToken, error) {
	exp := time.Now().Add(time.Second * time.Duration(ACCESS_TOKEN_EXPIRE))
	claims := jwt.MapClaims{
		"sub": user.Id.String(),
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(getJwtSecretKey())
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		UserId:    user.Id,
		Token:     signedStr,
		ExpiresIn: exp,
	}, nil
}

func (at* AccessToken) IsExpired() bool {
	return at.ExpiresIn.Before(time.Now())
}

type RefreshToken struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Token     string
	ExpiresIn time.Time
}

func NewRefreshToken(user *User) (*RefreshToken, error) {
	tokenId := uuid.New()
	exp := time.Now().Add(time.Second * time.Duration(REFRESH_TOKEN_EXPIRE))
	claims := jwt.MapClaims{
		"sub": user.Id.String(),
		"tid": tokenId.String(),
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(getJwtSecretKey())
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		Id:        tokenId,
		UserId:    user.Id,
		Token:     signedStr,
		ExpiresIn: exp,
	}, nil
}

func (rt* RefreshToken) IsExpired() bool {
	return rt.ExpiresIn.Before(time.Now())
}

// JWTの暗号化鍵
func getJwtSecretKey() []byte {
	secret, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		panic("\"JWT_SECRET_KEY\" is not set")
	}
	return []byte(secret)
}
