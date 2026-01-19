package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type GooglePsqlAndApiRepository struct {
	db           *sqlx.DB
	clientId     string
	clientSecret string
	redirectUri  string
}

func NewGooglePsqlAndApiRepository(
	db *sqlx.DB,
	clientId string,
	clientSecret string,
	redirectUri string,
) *GooglePsqlAndApiRepository {
	return &GooglePsqlAndApiRepository{
		db:           db,
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
	}
}

type GoogleOauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r *GooglePsqlAndApiRepository) GetTokenByCode(
	ctx context.Context,
	code string,
) (string, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", r.clientId)
	values.Add("client_secret", r.clientSecret)
	values.Add("redirect_uri", r.redirectUri)
	values.Add("grant_type", "authorization_code")
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", values)
	if err != nil {
		logrus.WithError(err).Error("failed to get oauth token")
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		logrus.WithField("statusCode", resp.StatusCode).Error("failed to get oauth token")
		return "", errors.New("token request status error")
	}
	defer resp.Body.Close()
	var response GoogleOauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.WithError(err).Error("token response decode error")
		return "", err
	}
	return response.AccessToken, nil
}

type GoogleUserInfoResponse struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func (r *GooglePsqlAndApiRepository) GetUserByToken(
	ctx context.Context,
	token string,
) (*user.GoogleUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("userinfo request error")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logrus.WithField("statusCode", resp.StatusCode).Error("userinfo response invalid status")
		return nil, err
	}
	defer resp.Body.Close()
	var response GoogleUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.WithError(err).Error("userinfo response decode error")
		return nil, err
	}
	googleUser, err := user.NewGoogleUser(response.Sub, response.Name, response.Email, response.Picture)
	if err != nil {
		return nil, err
	}
	return googleUser, nil
}

func (r *GooglePsqlAndApiRepository) CreateUser(
	ctx context.Context,
	user *user.User,
	avatar *user.Avatar,
	googleUser *user.GoogleUser,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO users(id, name, email, password_hash, created_at, updated_at)
			VALUES(:id, :name, :email, :password, :createdAt, :updatedAt)
		`
		_, err := tx.NamedExecContext(ctx, query, map[string]any{
			"id":        user.Id,
			"name":      user.Name,
			"email":     user.Email,
			"password":  user.Password,
			"createdAt": time.Now().UTC(),
			"updatedAt": time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		query = `
			INSERT INTO google_ids(user_id, google_id, created_at, updated_at)
			VALUES (:userId, :googleId, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"userId":    user.Id,
			"googleId":  googleUser.Id,
			"createdAt": time.Now().UTC(),
			"updatedAt": time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		if !avatar.Valid {
			return nil
		}
		query = `
			INSERT INTO avatars(id, url)
			VALUES(:id, :url)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":  avatar.Id,
			"url": avatar.Url,
		})
		if err != nil {
			return err
		}
		query = `
			INSERT INTO user_profiles(user_id, avatar_id)
			VALUES(:userId, :avatarId)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"userId":   user.Id,
			"avatarId": avatar.Id,
		})
		return err
	})
}
