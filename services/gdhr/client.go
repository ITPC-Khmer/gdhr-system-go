// Package gdhr is a thin client for the external GDHR public API. Every endpoint
// is a POST carrying the credentials in the JSON body; paged endpoints take the
// page number in the URL path.
package gdhr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client talks to the GDHR API.
type Client struct {
	baseURL string
	creds   credentials
	http    *http.Client
}

type credentials struct {
	U   string `json:"u"`
	P   string `json:"p"`
	Key string `json:"key"`
}

// New builds a GDHR client.
func New(baseURL, user, pass, key string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		creds:   credentials{U: user, P: pass, Key: key},
		http:    &http.Client{Timeout: timeout},
	}
}

// post sends the credentialed POST to path and decodes the JSON body into out.
func (c *Client) post(ctx context.Context, path string, out any) error {
	body, err := json.Marshal(c.creds)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gdhr %s -> HTTP %d: %s", path, resp.StatusCode, truncate(string(data), 300))
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("gdhr %s decode failed: %w", path, err)
	}
	return nil
}

// FetchInstitutes returns one page of institutes.
func (c *Client) FetchInstitutes(ctx context.Context, page int) (*InstitutesResponse, error) {
	var out InstitutesResponse
	if err := c.post(ctx, "/public-institutes/"+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchStaffs returns one page of staffs.
func (c *Client) FetchStaffs(ctx context.Context, page int) (*StaffsResponse, error) {
	var out StaffsResponse
	if err := c.post(ctx, "/public-staffs/"+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchRanks returns the full (unpaged) ranks list.
func (c *Client) FetchRanks(ctx context.Context) (*RanksResponse, error) {
	var out RanksResponse
	if err := c.post(ctx, "/public-ranks", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchPositions returns the full (unpaged) positions list.
func (c *Client) FetchPositions(ctx context.Context) (*PositionsResponse, error) {
	var out PositionsResponse
	if err := c.post(ctx, "/public-positions", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}
