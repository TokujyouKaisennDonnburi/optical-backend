package calendar

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// 招待状
type Invitation struct {
	Id           uuid.UUID
	CalendarId   uuid.UUID
	Email        string        // 招待先メアド(記録用)
	JoinedUserId uuid.NullUUID // 参加したユーザーのID(メアド違いで登録しても追跡するため)
	Token        uuid.UUID
	ExpiresAt    time.Time    // トークン有効期限
	UsedAt       sql.NullTime // 参加日時
	CreatedAt    time.Time
}

// 招待状を新規作成する(有効期限は30日)
func NewInvitation(calendarId uuid.UUID, email string) (*Invitation, error) {
	if calendarId == uuid.Nil {
		return nil, errors.New("Invitation `calendarId` is nil")
	}
	if email == "" {
		return nil, errors.New("Invitation `email` is empty")
	}
	return &Invitation{
		Id:         uuid.New(),
		CalendarId: calendarId,
		Email:      email,
		Token:      uuid.New(),
		ExpiresAt:  time.Now().UTC().AddDate(0, 0, 30),
		CreatedAt:  time.Now().UTC(),
	}, nil
}

// トークンが有効期限切れかどうか
func (i *Invitation) IsExpired() bool {
	return time.Now().UTC().After(i.ExpiresAt)
}

// 招待が使用済みかどうか
func (i *Invitation) IsUsed() bool {
	return i.UsedAt.Valid
}

// 招待を使用済みにする
func (i *Invitation) MarkAsUsed(userId uuid.UUID) {
	i.JoinedUserId = uuid.NullUUID{UUID: userId, Valid: true}
	i.UsedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
}
