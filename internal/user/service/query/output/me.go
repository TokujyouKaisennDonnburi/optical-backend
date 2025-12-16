package output

import (
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type UserQueryOutput struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Avatar    user.Avatar
	CreatedAt time.Time
	UpdatedAt time.Time
}
