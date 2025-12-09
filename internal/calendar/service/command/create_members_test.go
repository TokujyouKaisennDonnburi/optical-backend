package command

import (
	"context"
	"errors"
	"testing"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockMemberRepository は MemberRepository インターフェースのモック実装
type MockMemberRepository struct {
	CreateFunc func(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error
	JoinFunc   func(ctx context.Context, userId, calendarId uuid.UUID) error
}

func (m *MockMemberRepository) Create(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, userId, calendarId, emails)
	}
	return nil
}

func (m *MockMemberRepository) Join(ctx context.Context, userId, calendarId uuid.UUID) error {
	if m.JoinFunc != nil {
		return m.JoinFunc(ctx, userId, calendarId)
	}
	return nil
}

func TestCreateMember_Success(t *testing.T) {
	mockMemberRepo := &MockMemberRepository{
		CreateFunc: func(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error {
			return nil
		},
	}

	cmd := &CalendarCommand{
		memberRepository: mockMemberRepo,
	}

	input := MemberCreateInput{
		UserId:     uuid.New(),
		CalendarId: uuid.New(),
		Emails:     []string{"test@example.com", "user@domain.com"},
	}

	err := cmd.CreateMember(context.Background(), input)

	assert.NoError(t, err)
}

func TestCreateMember_InvalidEmail_Short(t *testing.T) {
	mockMemberRepo := &MockMemberRepository{}

	cmd := &CalendarCommand{
		memberRepository: mockMemberRepo,
	}

	input := MemberCreateInput{
		UserId:     uuid.New(),
		CalendarId: uuid.New(),
		Emails:     []string{"a@"}, // 2文字で短すぎる
	}

	err := cmd.CreateMember(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidEmail, err)
}

func TestCreateMember_InvalidEmail_NoAtSign(t *testing.T) {
	mockMemberRepo := &MockMemberRepository{}

	cmd := &CalendarCommand{
		memberRepository: mockMemberRepo,
	}

	input := MemberCreateInput{
		UserId:     uuid.New(),
		CalendarId: uuid.New(),
		Emails:     []string{"invalidemail.com"}, // @がない
	}

	err := cmd.CreateMember(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidEmail, err)
}

func TestCreateMember_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")
	mockMemberRepo := &MockMemberRepository{
		CreateFunc: func(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error {
			return expectedErr
		},
	}

	cmd := &CalendarCommand{
		memberRepository: mockMemberRepo,
	}

	input := MemberCreateInput{
		UserId:     uuid.New(),
		CalendarId: uuid.New(),
		Emails:     []string{"test@example.com"},
	}

	err := cmd.CreateMember(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCreateMember_MultipleEmails_OneInvalid(t *testing.T) {
	mockMemberRepo := &MockMemberRepository{}

	cmd := &CalendarCommand{
		memberRepository: mockMemberRepo,
	}

	input := MemberCreateInput{
		UserId:     uuid.New(),
		CalendarId: uuid.New(),
		Emails:     []string{"valid@example.com", "invalid"}, // 2番目が無効
	}

	err := cmd.CreateMember(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidEmail, err)
}

func TestCreateMember_EmptyEmails(t *testing.T) {
	called := false
	mockMemberRepo := &MockMemberRepository{
		CreateFunc: func(ctx context.Context, userId, calendarId uuid.UUID, emails []user.Email) error {
			called = true
			assert.Empty(t, emails)
			return nil
		},
	}

	cmd := &CalendarCommand{
		memberRepository: mockMemberRepo,
	}

	input := MemberCreateInput{
		UserId:     uuid.New(),
		CalendarId: uuid.New(),
		Emails:     []string{},
	}

	err := cmd.CreateMember(context.Background(), input)

	assert.NoError(t, err)
	assert.True(t, called)
}
