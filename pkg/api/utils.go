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
	"github.com/sirupsen/logrus"
)

const (
	GITHUB_API_VERSION = "2022-11-28"
	GITHUB_BASE_URL = "https://api.github.com"
)

func setRequestHeader(r *http.Request) error {
	token, err := getGithubAppBearerToken()
	if err != nil {
		return err
	}
	r.Header.Add("Accept", "application/vnd.github+json")
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("X-GitHub-Api-Version", GITHUB_API_VERSION)
	return nil
}

func getGithubAppBearerToken() (string, error) {
	filePath, ok := os.LookupEnv("GITHUB_PRIVATE_KEY_PEM_PATH")
	if !ok {
		return "", errors.New("'GITHUB_PRIVATE_KEY_PEM_PATH' is not set")
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		logrus.WithError(err).Error("failed to read github private key pem")
		return "", err
	}
	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logrus.Error("RSA private key block error")
		return "", errors.New("RSA private key block error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logrus.WithError(err).Error("failed to parse github private key pem")
		return "", err
	}
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
		"iss": auth.GetClientId(),
	})
	return token.SignedString(privateKey)
}
