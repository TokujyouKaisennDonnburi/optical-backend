package handler

import (
	"net/http"

	"github.com/google/uuid"
)

type GithubMilestonesRequest struct {
	UserId uuid.UUID
	ClendarId uuid.UUID
}

type GithubMilestonesResponse struct {
	Title string
	Progress int8
	Open int8
	Close int8
}

func(h *GithubHandler) GetMilestones (w http.ResponseWriter, r *http.Request){
	userId, err := uuid.Parse("")

}
