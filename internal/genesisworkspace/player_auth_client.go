package genesisworkspace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PlayerProfile struct {
	SubjectID   string `json:"subject_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type PlayerSession struct {
	Status    string        `json:"status"`
	Token     string        `json:"token,omitempty"`
	ExpiresAt string        `json:"expires_at,omitempty"`
	Player    PlayerProfile `json:"player"`
}

type PlayerRegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type PlayerLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PlayerAuthClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

type PlayerAuthHTTPError struct {
	Path       string
	StatusCode int
	Body       string
}

func (e *PlayerAuthHTTPError) Error() string {
	return fmt.Sprintf("IAOS Genesis authentication %s returned %d: %s", e.Path, e.StatusCode, e.Body)
}

func (e *PlayerAuthHTTPError) UpstreamStatusCode() int { return e.StatusCode }

func (c PlayerAuthClient) Register(ctx context.Context, input PlayerRegisterRequest) (PlayerSession, error) {
	var result PlayerSession
	err := c.doJSON(ctx, http.MethodPost, "/api/v1/genesis/auth/register", "", input, &result)
	return result, err
}

func (c PlayerAuthClient) Login(ctx context.Context, input PlayerLoginRequest) (PlayerSession, error) {
	var result PlayerSession
	err := c.doJSON(ctx, http.MethodPost, "/api/v1/genesis/auth/login", "", input, &result)
	return result, err
}

func (c PlayerAuthClient) Session(ctx context.Context, token string) (PlayerSession, error) {
	var result PlayerSession
	err := c.doJSON(ctx, http.MethodGet, "/api/v1/genesis/auth/session", token, nil, &result)
	return result, err
}

func (c PlayerAuthClient) doJSON(ctx context.Context, method, path, token string, input, output any) error {
	if strings.TrimSpace(c.BaseURL) == "" {
		return fmt.Errorf("IAOS base URL is required for Genesis authentication")
	}
	var body io.Reader
	if input != nil {
		raw, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	request, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.BaseURL, "/")+path, body)
	if err != nil {
		return err
	}
	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token = strings.TrimSpace(token); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &PlayerAuthHTTPError{
			Path: path, StatusCode: response.StatusCode, Body: strings.TrimSpace(string(raw)),
		}
	}
	if output != nil && len(raw) > 0 {
		return json.Unmarshal(raw, output)
	}
	return nil
}
