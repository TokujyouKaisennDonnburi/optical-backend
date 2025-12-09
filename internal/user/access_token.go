package user

import (
	"errors"
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

func NewAccessToken(userId uuid.UUID, userName string) (*AccessToken, error) {
	exp := time.Now().Add(time.Second * time.Duration(ACCESS_TOKEN_EXPIRE))
	claims := jwt.MapClaims{
		"sub":  userId.String(),
		"name": userName,
		"exp":  exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(getJwtSecretKey())
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		UserId:    userId,
		Token:     signedStr,
		ExpiresIn: exp,
	}, nil
}

func (at *AccessToken) IsExpired() bool {
	return at.ExpiresIn.Before(time.Now())
}

type RefreshToken struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	UserName  string
	Token     string
	ExpiresIn time.Time
}

func NewRefreshToken(user *User) (*RefreshToken, error) {
	tokenId := uuid.New()
	exp := time.Now().Add(time.Second * time.Duration(REFRESH_TOKEN_EXPIRE))
	claims := jwt.MapClaims{
		"sub":  user.Id.String(),
		"name": user.Name,
		"tid":  tokenId.String(),
		"exp":  exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(getJwtSecretKey())
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		Id:        tokenId,
		UserId:    user.Id,
		UserName:  user.Name,
		Token:     signedStr,
		ExpiresIn: exp,
	}, nil
}

func (rt *RefreshToken) IsExpired() bool {
	return rt.ExpiresIn.Before(time.Now())
}

// トークンをデコードして情報を取得する
func DecodeRefreshToken(refreshToken string) (*RefreshToken, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(getJwtSecretKey()), nil
	})
	if err != nil {
		return nil, err
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}
	name, ok := claims["name"]
	if !ok {
		return nil, errors.New("tid does not exist.")
	}
	userName, ok := name.(string)
	if !ok {
		return nil, errors.New("invalid userName.")
	}
	tid, ok := claims["tid"]
	if !ok {
		return nil, errors.New("tid does not exist.")
	}
	tidStr, ok := tid.(string)
	if !ok {
		return nil, errors.New("invalid tid.")
	}
	tidUuid, err := uuid.Parse(tidStr)
	if err != nil {
		return nil, err
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		Id:        tidUuid,
		UserId:    userId,
		UserName:  userName,
		Token:     refreshToken,
		ExpiresIn: exp.Time,
	}, nil
}

// JWTの暗号化鍵
func getJwtSecretKey() []byte {
	secret, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		panic("\"JWT_SECRET_KEY\" is not set")
	}
	return []byte(secret)
}
