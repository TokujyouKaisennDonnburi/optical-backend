package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type InstallationGetResponse struct {
	Id      int                            `json:"id"`
	Account InstallationGetResponseAccount `json:"account"`
}

type InstallationGetResponseAccount struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

func GetGithubInstallation(ctx context.Context, installationId string) (*InstallationGetResponse, error) {
	client := http.Client{}
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response InstallationGetResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
