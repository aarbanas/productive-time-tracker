package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.productive.io/api/v2"

// Client performs authenticated requests to the Productive API.
type Client struct {
	httpClient *http.Client
	token      string
	orgID      string
}

// NewClient returns a client configured with the given credentials.
func NewClient(token, orgID string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		orgID:      orgID,
	}
}

func (c *Client) get(endpoint string) ([]byte, error) {
	url := baseURL + endpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", c.token)
	req.Header.Set("X-Organization-Id", c.orgID)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) GetTimeReport(after, before string) ([]TimeReport, error) {
	url := fmt.Sprintf("/reports/time_reports?filter[after]=%s&filter[before]=%s", after, before)
	body, err := c.get(url)
	if err != nil {
		return nil, err
	}

	var response TimeReportsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse time reports response: %w", err)
	}

	return response.Data, nil
}
