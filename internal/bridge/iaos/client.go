// Package iaos implements the governed journal/cursor bridge; it never accesses IAOS DB or NATS.
package iaos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	base  *url.URL
	token string
	http  *http.Client
}

func New(baseURL, token string) (*Client, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("valid IAOS URL is required")
	}
	if token == "" {
		return nil, fmt.Errorf("short-lived token is required")
	}
	return &Client{u, token, &http.Client{Timeout: 15 * time.Second}}, nil
}

type Accept struct {
	MessageID     string `json:"message_id"`
	JournalCursor int64  `json:"journal_cursor"`
	Accepted      bool   `json:"accepted"`
	Duplicate     bool   `json:"duplicate"`
	RecordedAt    string `json:"recorded_at"`
	OperationRef  string `json:"operation_ref"`
}

func (c *Client) PostObservation(ctx context.Context, value worldcontract.Observation) (Accept, error) {
	return c.post(ctx, "/api/v1/world-bridge/observations", value)
}
func (c *Client) post(ctx context.Context, path string, value any) (Accept, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return Accept{}, err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.base.String()+path, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return Accept{}, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return Accept{}, fmt.Errorf("IAOS world bridge status %d: %s", res.StatusCode, string(body))
	}
	var out Accept
	if err := json.Unmarshal(body, &out); err != nil {
		return out, err
	}
	return out, nil
}

type Page struct {
	Items      []worldcontract.Envelope `json:"items"`
	NextCursor int64                    `json:"next_cursor"`
	HasMore    bool                     `json:"has_more"`
}

func (c *Client) Entries(ctx context.Context, runID, branch string, after int64) (Page, error) {
	q := url.Values{"world_run_id": {runID}, "branch_id": {branch}, "after": {strconv.FormatInt(after, 10)}, "limit": {"200"}}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.base.String()+"/api/v1/world-bridge/entries?"+q.Encode(), nil)
	req.Header.Set("Authorization", "Bearer "+c.token)
	res, err := c.http.Do(req)
	if err != nil {
		return Page{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return Page{}, fmt.Errorf("IAOS world bridge read status %d", res.StatusCode)
	}
	var out Page
	if err := json.NewDecoder(io.LimitReader(res.Body, 2<<20)).Decode(&out); err != nil {
		return out, err
	}
	return out, nil
}
