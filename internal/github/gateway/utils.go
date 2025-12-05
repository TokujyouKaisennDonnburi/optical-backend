package gateway

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	GITHUB_PRIVATE_KEY_PEM_PATH = "optical-github.2025-12-04.private-key.pem"
)

func setRequestHeader(r *http.Request) error {
	token, err := getGithubAppBearerToken()
	if err != nil {
		return err
	}
	r.Header.Add("Accept", "application/vnd.github+json")
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	return nil
}

func getGithubAppBearerToken() (string, error) {
	file, err := os.ReadFile(GITHUB_PRIVATE_KEY_PEM_PATH)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("RSA private key block error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": getGithubAppId(),
	})
	return token.SignedString(privateKey)
}

func getGithubAppId() string {
	appId, ok := os.LookupEnv("GITHUB_APP_ID")
	if !ok {
		panic("'GITHUB_APP_ID' is not set")
	}
	return appId
}
