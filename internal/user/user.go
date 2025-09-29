package user

import (
	"errors"
	"time"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

var (
	ErrInvalidEmail    = errors.New("無効なメールアドレスです")
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
	return &User{
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
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
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// deleted chek
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// hash pass
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verify pass
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// validate name
func validateName(name string) error {
	if len(name) < 1 || len(name) > 50 {
		return ErrInvalidName
	}
	return nil
}

// validate email
func validateEmail(email string) error {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}
	return nil
}
// validate password
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	return nil
}

