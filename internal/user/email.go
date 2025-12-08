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
	length := utf8.RuneCountInString(email)
	if length < MIN_EMAIL_LENGTH || !strings.Contains(email, "@") {
		return "", ErrInvalidEmail
	}
	return Email(email), nil
}

func NewEmails(emails []string) ([]string, error){
	result := make([]string, 0, len(emails))
	for _, email := range emails {
		validated, err := NewEmail(email)
		if err != nil {
			return nil, err
		}
		result = append(result, string(validated))
	}
	return result, nil
}
