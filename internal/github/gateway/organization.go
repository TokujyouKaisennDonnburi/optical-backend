package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/security"
	"github.com/sirupsen/logrus"
)

// トークンを復号化して取得
func (r *StateRedisRepository) GetOrganization(
	ctx context.Context,
	installationId string,
) (*github.Organization, error) {
	// キャッシュから取得する関数を定義
	getFromCache := func() (*github.Organization, error) {
		redisKey := security.Hash(installationId, getEncryptionKey())
		result, err := r.redisClient.Get(ctx, getGithubAccessKey(redisKey)).Result()
		if err != nil {
			return nil, err
		}
		plainText, err := security.Decrypt(result, getEncryptionKey())
		if err != nil {
			return nil, err
		}
		var organization github.Organization
		if err := json.Unmarshal([]byte(plainText), &organization); err != nil {
			return nil, err
		}
		return &organization, nil
	}
	// キャッシュから取得する処理を実行
	organization, err := getFromCache()
	if err == nil {
		return organization, nil
	}
	logrus.WithError(err).Error("Cache get error: ")
	// 失敗した場合APIから取得
	organization, err = api.GetInstalledOrganization(ctx, installationId)
	if err != nil {
		return nil, err
	}
	// キャッシュ処理を行い失敗した場合エラーのみ出力
	cacheFunc := func() error {
		now := time.Now().UTC()
		if organization.TokenExpiresAt.Before(now) {
			return errors.New("cache error: invalid time")
		}
		return r.SaveOrganization(ctx, organization, organization.TokenExpiresAt.Sub(now))
	}
	if err := cacheFunc(); err != nil {
		logrus.WithError(err).Error("cache error")
	}
	return organization, nil
}

// 暗号化してトークンを保存
func (r *StateRedisRepository) SaveOrganization(
	ctx context.Context,
	organization *github.Organization,
	expiration time.Duration,
) error {
	jsonText, err := json.Marshal(organization)
	if err != nil {
		return err
	}
	encryptedToken, err := security.Encrypt(string(jsonText), getEncryptionKey())
	if err != nil {
		return err
	}
	redisKey := security.Hash(organization.InstallationId, getEncryptionKey())
	return r.redisClient.SetEx(ctx, getGithubAccessKey(redisKey), encryptedToken, expiration).Err()
}

// TODO: DIにする
func getEncryptionKey() []byte {
	encryptionKey, ok := os.LookupEnv("REDIS_ENCRYPTION_KEY")
	if !ok {
		panic("'REDIS_ENCRYPTION_KEY' is not set")
	}
	return []byte(encryptionKey)
}

func getGithubAccessKey(installationId string) string {
	return "github:accessToken:" + installationId
}
