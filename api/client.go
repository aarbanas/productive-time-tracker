package api

import (
	"context"
	"fmt"
	"time"

	"github.com/aarbanas/productive-time-tracker/api/generated"
)

const serverURL = "https://api.productive.io/api/v2"

type tokenSource struct {
	token string
}

func (t *tokenSource) HeaderToken(_ context.Context, _ generated.OperationName) (generated.HeaderToken, error) {
	var h generated.HeaderToken
	h.SetAPIKey(t.token)
	return h, nil
}

// ProductiveClient wraps the generated ogen client with stored credentials.
type ProductiveClient struct {
	inner *generated.Client
	orgID string
}

// NewProductiveClient returns a client authenticated with the given token and org ID.
func NewProductiveClient(token, orgID string) (*ProductiveClient, error) {
	c, err := generated.NewClient(serverURL, &tokenSource{token: token})
	if err != nil {
		return nil, err
	}
	return &ProductiveClient{inner: c, orgID: orgID}, nil
}

// GetTimeReport fetches aggregated time report data for the given date range (YYYY-MM-DD).
func (c *ProductiveClient) GetTimeReport(after, before string) ([]generated.TimeReport, error) {
	afterDate, err := time.Parse("2006-01-02", after)
	if err != nil {
		return nil, fmt.Errorf("invalid after date: %w", err)
	}
	beforeDate, err := time.Parse("2006-01-02", before)
	if err != nil {
		return nil, fmt.Errorf("invalid before date: %w", err)
	}

	resp, err := c.inner.ReportsTimeReportsIndex(context.Background(), generated.ReportsTimeReportsIndexParams{
		XOrganizationID: c.orgID,
		FilterAfter:     generated.NewOptDate(afterDate),
		FilterBefore:    generated.NewOptDate(beforeDate),
	})
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
