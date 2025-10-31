package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query"
)

type UserHttpHandler struct {
	userQuery   *query.UserQuery
	userCommand *command.UserCommand
}

func NewUserHttpHandler(userQuery *query.UserQuery, userCommand *command.UserCommand) *UserHttpHandler {
	if userQuery == nil {
		panic("UserQuery is nil")
	}
	if userCommand == nil {
		panic("UserCommand is nil")
	}
	return &UserHttpHandler{
		userQuery:   userQuery,
		userCommand: userCommand,
	}
}
