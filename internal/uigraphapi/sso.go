package uigraphapi

import (
	"context"
	"time"
)

type OAuthProvider struct {
	ID             string    `json:"id"`
	ProviderName   string    `json:"providerName"`
	Type           string    `json:"type"`
	DisplayName    string    `json:"displayName"`
	IconURL        string    `json:"iconUrl"`
	ClientID       string    `json:"clientId"`
	ClientSecret   string    `json:"clientSecret"`
	AuthURL        string    `json:"authUrl"`
	TokenURL       string    `json:"tokenUrl"`
	UserinfoURL    string    `json:"userinfoUrl"`
	APIURL         string    `json:"apiUrl"`
	Scopes         string    `json:"scopes"`
	AllowedDomains string    `json:"allowedDomains"`
	AllowSignUp    bool      `json:"allowSignUp"`
	EmailClaim     string    `json:"emailClaim"`
	NameClaim      string    `json:"nameClaim"`
	SubClaim       string    `json:"subClaim"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type RoleMapping struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	ClaimKey       string `json:"claimKey"`
	ClaimValue     string `json:"claimValue"`
	Role           string `json:"role"`
	Scope          string `json:"scope"`
	ResourceType   string `json:"resourceType"`
	ResourceID     string `json:"resourceId"`
}

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

type SAMLConfig struct {
	ID              string    `json:"id"`
	IDPMetadataURL  string    `json:"idpMetadataUrl"`
	IDPMetadataXML  string    `json:"idpMetadataXml"`
	IDPEntityID     string    `json:"idpEntityId"`
	IDPSsoURL       string    `json:"idpSsoUrl"`
	IDPCert         string    `json:"idpCert"`
	SPEntityID      string    `json:"spEntityId"`
	SPCert          string    `json:"spCert"`
	SPKey           string    `json:"spKey"`
	SignRequests    bool      `json:"signRequests"`
	NameIDFormat    string    `json:"nameIdFormat"`
	EmailAttribute  string    `json:"emailAttribute"`
	NameAttribute   string    `json:"nameAttribute"`
	LoginAttribute  string    `json:"loginAttribute"`
	GroupsAttribute string    `json:"groupsAttribute"`
	AllowSignUp     bool      `json:"allowSignUp"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

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

func (c *Client) PrepareOAuthProviderIconUpload(ctx context.Context, provider string) (*AssetUpload, error) {
	var out AssetUpload
	return &out, c.post(ctx, "/api/v1/sso/oauth/"+provider+"/icon/prepare", nil, &out)
}

func (c *Client) SetOAuthProviderIcon(ctx context.Context, provider string) error {
	return c.put(ctx, "/api/v1/sso/oauth/"+provider+"/icon", nil, nil)
}

func (c *Client) RemoveOAuthProviderIcon(ctx context.Context, provider string) error {
	return c.del(ctx, "/api/v1/sso/oauth/"+provider+"/icon")
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

// ── SCIM ──────────────────────────────────────────────────────────────────────

type SCIMConfig struct {
	ID string `json:"id"`
}

func (c *Client) GetSCIM(ctx context.Context) (*SCIMConfig, error) {
	var out SCIMConfig
	return &out, c.get(ctx, "/api/v1/sso/scim", &out)
}
