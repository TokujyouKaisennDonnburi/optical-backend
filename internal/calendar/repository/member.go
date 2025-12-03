package repository

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"

type MemberRepository interface {
	Create(Email string)(calendar.Member, error)
}
