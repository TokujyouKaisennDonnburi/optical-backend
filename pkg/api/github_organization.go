package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
)

type GithubAppAccessTokenResponse struct {
	Token        string                     `json:"token"`
	ExpiresAt    time.Time                  `json:"expires_at"`
	Repositories []GithubRepositoryResponse `json:"repositories"`
}

func GetInstalledOrganization(ctx context.Context, installationId string) (*github.Organization, error) {
	client := http.Client{}
	// インストール済みアプリ取得リクエスト
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		GITHUB_BASE_URL+"/app/installations/"+installationId,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = setRequestHeader(req)
	if err != nil {
		return nil, err
	}
	respGet, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if respGet.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get github app installation: %d", respGet.StatusCode)
	}
	defer respGet.Body.Close()
	var respGetBody InstallationGetResponse
	if err = json.NewDecoder(respGet.Body).Decode(&respGetBody); err != nil {
		return nil, err
	}
	// アクセストークン発行リクエスト
	req, err = http.NewRequest(
		"POST",
		GITHUB_BASE_URL+"/app/installations/"+installationId+"/access_tokens",
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = setRequestHeader(req)
	if err != nil {
		return nil, err
	}
	respPost, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if respPost.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create access token: %d", respPost.StatusCode)
	}
	defer respPost.Body.Close()
	var respPostBody GithubAppAccessTokenResponse
	if err := json.NewDecoder(respPost.Body).Decode(&respPostBody); err != nil {
		return nil, err
	}
	repositories, err := GetInstalledRepositories(ctx, respPostBody.Token)
	if err != nil {
		return nil, err
	}
	return &github.Organization{
		Id:             respGetBody.Account.Id,
		Login:          respGetBody.Account.Login,
		InstallationId: installationId,
		AccessToken:    respPostBody.Token,
		TokenExpiresAt: respPostBody.ExpiresAt,
		Repositories:   repositories,
	}, nil
}

