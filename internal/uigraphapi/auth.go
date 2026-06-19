package uigraphapi

import "context"

type MeResponse struct {
	UserID       string `json:"userId"`
	OrgID        string `json:"orgId"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Login        string `json:"login"`
	Kind         string `json:"kind"`
	Role         string `json:"role"`
	AuthProvider string `json:"authProvider"`
	AvatarURL    string `json:"avatarUrl,omitempty"`
}

type OrgSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

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
