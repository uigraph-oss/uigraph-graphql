package uigraphapi

import (
	"context"
	"time"
)

type LDAPConfig struct {
	ID                string    `json:"id"`
	Host              string    `json:"host"`
	Port              int       `json:"port"`
	UseSSL            bool      `json:"useSsl"`
	StartTLS          bool      `json:"startTls"`
	SkipTLSVerify     bool      `json:"skipTlsVerify"`
	BindDN            string    `json:"bindDn"`
	BindPassword      string    `json:"bindPassword"`
	SearchBaseDN      string    `json:"searchBaseDn"`
	SearchFilter      string    `json:"searchFilter"`
	EmailAttribute    string    `json:"emailAttribute"`
	NameAttribute     string    `json:"nameAttribute"`
	UsernameAttribute string    `json:"usernameAttribute"`
	MemberOfAttribute string    `json:"memberOfAttribute"`
	AllowSignUp       bool      `json:"allowSignUp"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

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

type SCIMConfig struct {
	ID string `json:"id"`
}

func (c *Client) GetSCIM(ctx context.Context) (*SCIMConfig, error) {
	var out SCIMConfig
	return &out, c.get(ctx, "/api/v1/sso/scim", &out)
}
