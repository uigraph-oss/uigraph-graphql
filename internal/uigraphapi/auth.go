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
	ID      string `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl,omitempty"`
	Role           string `json:"role"`
	Active         bool   `json:"active"`
	OnboardingDone bool   `json:"onboardingDone"`
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

func (c *Client) PrepareUserAvatarUpload(ctx context.Context) (*AssetUpload, error) {
	var out AssetUpload
	return &out, c.post(ctx, "/api/v1/users/me/avatar/prepare", nil, &out)
}

func (c *Client) SetMyAvatar(ctx context.Context) error {
	return c.put(ctx, "/api/v1/users/me/avatar", nil, nil)
}
