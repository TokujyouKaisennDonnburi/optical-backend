package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
)

type InstallationRepositoryResponse struct {
	TotalCount   int                        `json:"total_count"`
	Repositories []GithubRepositoryResponse `json:"repositories"`
}

type GithubRepositoryResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func GetInstalledRepositories(ctx context.Context, accessToken string) ([]github.Repository, error) {
	client := http.Client{}
	// インストール済みアプリ取得リクエスト
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		GITHUB_BASE_URL+"/installation/repositories",
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = setRequestHeader(req)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	respGet, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if respGet.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repositories: %d", respGet.StatusCode)
	}
	defer respGet.Body.Close()
	var respGetBody InstallationRepositoryResponse
	if err = json.NewDecoder(respGet.Body).Decode(&respGetBody); err != nil {
		return nil, err
	}
	var repositories []github.Repository
	for _, repository := range respGetBody.Repositories {
		repositories = append(repositories, github.Repository{
			Id:       repository.Id,
			Name:     repository.Name,
			FullName: repository.FullName,
		})
	}
	fmt.Println("Repositories", repositories)
	return repositories, nil
}
