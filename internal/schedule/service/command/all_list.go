package command

import (
		"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
		"github.com/google/uuid"
)

type SchduleAllList struct{
	Id uuid.UUID
	Name string
}
