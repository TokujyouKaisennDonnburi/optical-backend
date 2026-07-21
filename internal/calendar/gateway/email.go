package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
)

// Resend APIのエンドポイント
const resendApiUrl = "https://api.resend.com/emails"

// Resendは1リクエストあたり最大50宛先までしか送れない。
// toに自分のアドレスを1件使うため、bccは49件ずつに分割して複数リクエストに分ける。
const maxBccPerRequest = 49

type ResendRepository struct {
	apiKey     string
	fromEmail  string
	httpClient *http.Client
}

func NewResendRepository(apiKey string, fromEmail string) *ResendRepository {
	if apiKey == "" {
		panic("resend api key is empty")
	}
	if fromEmail == "" {
		panic("email address is empty")
	}
	return &ResendRepository{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// resendEmailRequest はResend APIへのリクエストボディ
type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

func (r *ResendRepository) NotifyAll(ctx context.Context, subject, content string, emails []user.Email) error {
	if len(emails) == 0 {
		return nil
	}

	mails := make([]string, len(emails))
	for i, email := range emails {
		mails[i] = string(email)
	}

	// bccを49件ずつに分割して複数リクエストに分ける。
	// 途中のリクエストが失敗した場合、それ以前のリクエスト分は送信済みとなる点に注意。
	for i := 0; i < len(mails); i += maxBccPerRequest {
		end := i + maxBccPerRequest
		if end > len(mails) {
			end = len(mails)
		}
		if err := r.send(ctx, subject, content, mails[i:end]); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResendRepository) send(ctx context.Context, subject, content string, bcc []string) error {
	reqBody := resendEmailRequest{
		From:    r.fromEmail,
		To:      []string{r.fromEmail},
		Bcc:     bcc,
		Subject: subject,
		Text:    content,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendApiUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend api request failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}
