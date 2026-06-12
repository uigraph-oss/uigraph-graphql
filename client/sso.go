package client

import "context"

// ── OAuth providers ───────────────────────────────────────────────────────────

func (c *Client) ListOAuthProviders(ctx context.Context) ([]OAuthProvider, error) {
	var out struct {
		Providers []OAuthProvider `json:"providers"`
	}
	return out.Providers, c.get(ctx, "/api/v1/sso/oauth", &out)
}

func (c *Client) UpsertOAuthProvider(ctx context.Context, provider string, body map[string]interface{}) error {
	return c.put(ctx, "/api/v1/sso/oauth/"+provider, body, nil)
}

func (c *Client) DeleteOAuthProvider(ctx context.Context, provider string) error {
	return c.del(ctx, "/api/v1/sso/oauth/"+provider)
}

// ── Role mappings ─────────────────────────────────────────────────────────────

func (c *Client) ListRoleMappings(ctx context.Context) ([]RoleMapping, error) {
	var out struct {
		Mappings []RoleMapping `json:"mappings"`
	}
	return out.Mappings, c.get(ctx, "/api/v1/sso/role-mappings", &out)
}

func (c *Client) CreateRoleMapping(ctx context.Context, body map[string]interface{}) error {
	return c.post(ctx, "/api/v1/sso/role-mappings", body, nil)
}

func (c *Client) DeleteRoleMapping(ctx context.Context, id string) error {
	return c.del(ctx, "/api/v1/sso/role-mappings/"+id)
}

// ── LDAP ──────────────────────────────────────────────────────────────────────

func (c *Client) GetLDAP(ctx context.Context) (*LDAPConfig, error) {
	var out LDAPConfig
	return &out, c.get(ctx, "/api/v1/sso/ldap", &out)
}

func (c *Client) UpsertLDAP(ctx context.Context, body map[string]interface{}) error {
	return c.put(ctx, "/api/v1/sso/ldap", body, nil)
}

func (c *Client) DeleteLDAP(ctx context.Context) error {
	return c.del(ctx, "/api/v1/sso/ldap")
}

// ── SAML ──────────────────────────────────────────────────────────────────────

func (c *Client) GetSAML(ctx context.Context) (*SAMLConfig, error) {
	var out SAMLConfig
	return &out, c.get(ctx, "/api/v1/sso/saml", &out)
}

func (c *Client) UpsertSAML(ctx context.Context, body map[string]interface{}) error {
	return c.put(ctx, "/api/v1/sso/saml", body, nil)
}
