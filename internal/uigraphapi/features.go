package uigraphapi

import "context"

type Features struct {
	GitHub     bool   `json:"githubEnabled"`
	Enterprise bool   `json:"enterpriseEnabled"`
	BillingURL string `json:"billingUrl"`
}

func (c *Client) GetFeatures(ctx context.Context) (*Features, error) {
	var out Features
	if err := c.get(ctx, "/api/v1/instance-info", &out); err != nil {
		return nil, err
	}
	return &out, nil
}
