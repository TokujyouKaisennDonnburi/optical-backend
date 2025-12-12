package api

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
)

const (
	GITHUB_PRIVATE_KEY_PEM_PATH = "optical-github.private-key.pem"
	GITHUB_API_VERSION = "2022-11-28"
	GITHUB_BASE_URL = "https://api.github.com"
)

func SetRequestHeader(r *http.Request) error {
	token, err := GetGithubAppBearerToken()
	if err != nil {
		return err
	}
	r.Header.Add("Accept", "application/vnd.github+json")
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("X-GitHub-Api-Version", GITHUB_API_VERSION)
	return nil
}

func GetGithubAppBearerToken() (string, error) {
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
		"iss": auth.GetGithubAppId(),
	})
	return token.SignedString(privateKey)
}
