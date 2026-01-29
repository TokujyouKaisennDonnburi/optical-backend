package command

import (
	calendarRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/transact"
)

type UserCommand struct {
	transactor           transact.TransactionProvider
	userRepository       repository.UserRepository
	tokenRepository      repository.TokenRepository
	avatarRepository     repository.AvatarRepository
	googleRepository     repository.GoogleRepository
	calendarRepository   calendarRepo.CalendarRepository
	oauthStateRepository repository.OauthStateRepository
}

func NewUserCommand(
	transactor transact.TransactionProvider,
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	avatarRepository repository.AvatarRepository,
	googleRepository repository.GoogleRepository,
	calendarRepository calendarRepo.CalendarRepository,
	oauthStateRepository repository.OauthStateRepository,
) *UserCommand {
	if transactor == nil {
		panic("transactor is nil")
	}
	if userRepository == nil {
		panic("userRepository is nil")
	}
	if tokenRepository == nil {
		panic("tokenRepository is nil")
	}
	if avatarRepository == nil {
		panic("avatarRepository is nil")
	}
	if googleRepository == nil {
		panic("googleRepository is nil")
	}
	if calendarRepository == nil {
		panic("caledarRepository is nil")
	}
	if oauthStateRepository == nil {
		panic("oauthStateRepository is nil")
	}
	return &UserCommand{
		transactor:           transactor,
		userRepository:       userRepository,
		tokenRepository:      tokenRepository,
		avatarRepository:     avatarRepository,
		googleRepository:     googleRepository,
		calendarRepository:   calendarRepository,
		oauthStateRepository: oauthStateRepository,
	}
}
