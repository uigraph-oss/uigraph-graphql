package uigraphapi

import "context"

func (c *Client) Me(ctx context.Context) (*MeResponse, error) {
	var out MeResponse
	return &out, c.get(ctx, "/api/v1/auth/me", &out)
}

func (c *Client) MyOrgs(ctx context.Context) ([]OrgSummary, error) {
	var out struct {
		Orgs []OrgSummary `json:"orgs"`
	}
	return out.Orgs, c.get(ctx, "/api/v1/auth/orgs", &out)
}

func (c *Client) SwitchOrg(ctx context.Context, orgID string) error {
	return c.post(ctx, "/api/v1/auth/switch-org", map[string]string{"orgId": orgID}, nil)
}
