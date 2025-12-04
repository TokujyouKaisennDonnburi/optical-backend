package user

import (
	"errors"
	"strings"
	"unicode/utf8"
)
const MIN_EMAIL_LENGTH = 3

type Email string

var ErrInvalidEmail = errors.New("無効なメールアドレスです")

func NewEmail(email string)(Email,error){
	if email == "" {
		return "", ErrInvalidEmail
	}
	length := utf8.RuneCountInString(email)
	if length < MIN_EMAIL_LENGTH || !strings.Contains(email, "@") {
		return "", ErrInvalidEmail
	}
	return Email(email), nil
}

