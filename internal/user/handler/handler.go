package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/command"

type UserHttpHandler struct {
	userCommand *command.UserCommand
}

func NewUserHttpHandler(userCommand *command.UserCommand) *UserHttpHandler {
	if userCommand == nil {
		panic("UserCommand is nil")
	}
	return &UserHttpHandler{
		userCommand: userCommand,
	}
}
