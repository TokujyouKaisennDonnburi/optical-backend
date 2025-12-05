package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	GITHUB_BASE_URL = "https://api.github.com"
)

type GithubApiRepository struct{}

func NewGithubApiRepository() *GithubApiRepository {
	return &GithubApiRepository{}
}

type InstallationGetResponse struct {
	Id      int                            `json:"id"`
	Account InstallationGetResponseAccount `json:"account"`
}

type InstallationGetResponseAccount struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

func (r *GithubApiRepository) LinkUser(
	ctx context.Context,
	userId uuid.UUID,
	installationId string,
) error {
	client := http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		GITHUB_BASE_URL+"/app/installations/"+installationId,
		nil,
	)
	if err != nil {
		return err
	}
	err = setRequestHeader(req)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	var response InstallationGetResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Status: %d", resp.StatusCode)
	fmt.Printf("Body: %s", string(body))
	fmt.Printf("Github Id: %d\n", response.Account.Id)
	fmt.Printf("Github Login: %s\n", response.Account.Login)
	return nil
}
