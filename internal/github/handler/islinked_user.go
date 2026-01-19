package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

// 未連携時用でomitemptyを使用
type IsLinkedUserResponse struct {
	IsLinked    bool   `json:"isLinked"`
	GithubId    string `json:"githubId,omitempty"`
	GithubName  string `json:"githubName,omitempty"`
	GithubEmail string `json:"githubEmail,omitempty"`
	IsSsoLogin  bool   `json:"isSsoLogin,omitempty"`
	LinkedAt    string `json:"linkedAt,omitempty"`
}

func (h *GithubHandler) IsLinkedUser(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// handlerを使ってquery呼び出し
	result, err := h.githubQuery.IsLinkedUser(r.Context(), userId)
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	// linledat(time.time)が無い場合、空文字にして返却したい
	var linkedAt string
	if !result.LinkedAt.IsZero() {
		linkedAt = result.LinkedAt.UTC().Format(time.RFC3339)
	}

	render.JSON(w, r, IsLinkedUserResponse{
		IsLinked:    result.IsLinked,
		GithubId:    result.GithubId,
		GithubName:  result.GithubName,
		GithubEmail: result.GithubEmail,
		IsSsoLogin:  result.IsSsoLogin,
		LinkedAt:    linkedAt,
	})
}
