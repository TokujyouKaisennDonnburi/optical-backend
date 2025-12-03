package repository

import "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"

type MemberRepository interface {
	create(Email string)(calendar.Member, error)
}
