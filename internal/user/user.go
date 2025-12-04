package user

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	MIN_USER_NAME_LENGTH = 1
	MAX_USER_NAME_LENGTH = 50
	MIN_PASSWORD_LENGTH = 8
)

type User struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

var (
	ErrInvalidName     = errors.New("ユーザー名は1文字以上50文字以下である必要があります")
	ErrInvalidPassword = errors.New("パスワードは8文字以上である必要があります")
)

func NewUser(name, email, password string) (*User, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	// hashed pass
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &User{
		Id:		   id,
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: time.Time{},
	}, nil
}

// update name
func (u *User) UpdateName(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// update email
func (u *User) UpdateEmail(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// update pass
func (u *User) UpdatePassword(newPassword string) error {
	if err := validatePassword(newPassword); err != nil {
		return err
	}
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// delete user
func (u *User) Delete() {
	now := time.Now()
	u.DeletedAt = now
	u.UpdatedAt = now
}

// deleted check
func (u *User) IsDeleted() bool {
	return !u.DeletedAt.IsZero()
}

// hash pass
func hashPassword(password string) ([]byte, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedBytes, nil
}

// verify pass
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// validate name
func validateName(name string) error {
	length := utf8.RuneCountInString(name)
	if length < MIN_USER_NAME_LENGTH || length > MAX_USER_NAME_LENGTH {
		return ErrInvalidName
	}
	return nil
}

// validate pass
func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < MIN_PASSWORD_LENGTH {
		return ErrInvalidPassword
	}
	return nil
}
