package uigraphapi

import (
	"context"
	"time"
)

type AuthProvider struct {
	ID             string `json:"id"`
	Slug           string `json:"slug"`
	OrgID          string `json:"orgId"`
	Kind           string `json:"kind"`
	Type           string `json:"type"`
	DisplayName    string `json:"displayName"`
	IconURL        string `json:"iconUrl"`
	Enabled        bool   `json:"enabled"`
	AllowSignUp    bool   `json:"allowSignUp"`
	AllowedDomains string `json:"allowedDomains"`
	DefaultRole    string `json:"defaultRole"`

	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AuthURL      string `json:"authUrl"`
	TokenURL     string `json:"tokenUrl"`
	UserinfoURL  string `json:"userinfoUrl"`
	APIURL       string `json:"apiUrl"`
	Scopes       string `json:"scopes"`
	EmailClaim   string `json:"emailClaim"`
	NameClaim    string `json:"nameClaim"`
	SubClaim     string `json:"subClaim"`
	GroupsClaim  string `json:"groupsClaim"`

	IDPMetadataURL  string `json:"idpMetadataUrl"`
	IDPMetadataXML  string `json:"idpMetadataXml"`
	IDPEntityID     string `json:"idpEntityId"`
	IDPSsoURL       string `json:"idpSsoUrl"`
	IDPCert         string `json:"idpCert"`
	SPEntityID      string `json:"spEntityId"`
	SPCert          string `json:"spCert"`
	SPKey           string `json:"spKey"`
	SignRequests    bool   `json:"signRequests"`
	NameIDFormat    string `json:"nameIdFormat"`
	EmailAttribute  string `json:"emailAttribute"`
	NameAttribute   string `json:"nameAttribute"`
	GroupsAttribute string `json:"groupsAttribute"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SamlSpMetadata struct {
	EntityID    string `json:"entityId"`
	ACSURL      string `json:"acsUrl"`
	MetadataURL string `json:"metadataUrl"`
	Metadata    string `json:"metadata"`
}

type AuthRoleMapping struct {
	ID             string `json:"id"`
	ProviderID     string `json:"providerId"`
	Priority       int    `json:"priority"`
	AttributeKey   string `json:"attributeKey"`
	Operator       string `json:"operator"`
	AttributeValue string `json:"attributeValue"`
	Role           string `json:"role"`
}

type OrgDomain struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"orgId"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MappingOperator struct {
	Name       string `json:"name"`
	TakesValue bool   `json:"takesValue"`
}

func authProvidersPath(orgID string) string {
	return "/api/v1/orgs/" + orgID + "/auth/providers"
}

func (c *Client) ListAuthProviders(ctx context.Context, orgID string) ([]AuthProvider, error) {
	var out struct {
		Providers []AuthProvider `json:"providers"`
	}
	return out.Providers, c.get(ctx, authProvidersPath(orgID), &out)
}

func (c *Client) GetAuthProvider(ctx context.Context, orgID, slug string) (*AuthProvider, error) {
	var out AuthProvider
	return &out, c.get(ctx, authProvidersPath(orgID)+"/"+slug, &out)
}

func (c *Client) CreateAuthProvider(ctx context.Context, orgID string, body map[string]interface{}) (*AuthProvider, error) {
	var out AuthProvider
	return &out, c.post(ctx, authProvidersPath(orgID), body, &out)
}

func (c *Client) UpdateAuthProvider(ctx context.Context, orgID, slug string, body map[string]interface{}) (*AuthProvider, error) {
	var out AuthProvider
	return &out, c.put(ctx, authProvidersPath(orgID)+"/"+slug, body, &out)
}

func (c *Client) DeleteAuthProvider(ctx context.Context, orgID, slug string) error {
	return c.del(ctx, authProvidersPath(orgID)+"/"+slug)
}

func (c *Client) PrepareAuthProviderIconUpload(ctx context.Context, orgID, slug string) (*AssetUpload, error) {
	var out AssetUpload
	return &out, c.post(ctx, authProvidersPath(orgID)+"/"+slug+"/icon/prepare", nil, &out)
}

func (c *Client) SetAuthProviderIcon(ctx context.Context, orgID, slug string) error {
	return c.put(ctx, authProvidersPath(orgID)+"/"+slug+"/icon", nil, nil)
}

func (c *Client) RemoveAuthProviderIcon(ctx context.Context, orgID, slug string) error {
	return c.del(ctx, authProvidersPath(orgID)+"/"+slug+"/icon")
}

func (c *Client) GetAuthProviderSAMLMetadata(ctx context.Context, orgID, slug string) (*SamlSpMetadata, error) {
	var out SamlSpMetadata
	return &out, c.get(ctx, authProvidersPath(orgID)+"/"+slug+"/saml/metadata", &out)
}

func (c *Client) ListAuthRoleMappings(ctx context.Context, orgID, providerSlug string) ([]AuthRoleMapping, error) {
	var out struct {
		Mappings []AuthRoleMapping `json:"mappings"`
	}
	return out.Mappings, c.get(ctx, authProvidersPath(orgID)+"/"+providerSlug+"/role-mappings", &out)
}

func (c *Client) CreateAuthRoleMapping(ctx context.Context, orgID, providerSlug string, body map[string]interface{}) (*AuthRoleMapping, error) {
	var out AuthRoleMapping
	return &out, c.post(ctx, authProvidersPath(orgID)+"/"+providerSlug+"/role-mappings", body, &out)
}

func (c *Client) UpdateAuthRoleMapping(ctx context.Context, orgID, providerSlug, mappingID string, body map[string]interface{}) (*AuthRoleMapping, error) {
	var out AuthRoleMapping
	return &out, c.put(ctx, authProvidersPath(orgID)+"/"+providerSlug+"/role-mappings/"+mappingID, body, &out)
}

func (c *Client) DeleteAuthRoleMapping(ctx context.Context, orgID, providerSlug, mappingID string) error {
	return c.del(ctx, authProvidersPath(orgID)+"/"+providerSlug+"/role-mappings/"+mappingID)
}

func (c *Client) ListOrgDomains(ctx context.Context, orgID string) ([]OrgDomain, error) {
	var out struct {
		Domains []OrgDomain `json:"domains"`
	}
	return out.Domains, c.get(ctx, "/api/v1/orgs/"+orgID+"/auth/domains", &out)
}

func (c *Client) CreateOrgDomain(ctx context.Context, orgID, domain string) (*OrgDomain, error) {
	var out OrgDomain
	body := map[string]interface{}{"domain": domain}
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/auth/domains", body, &out)
}

func (c *Client) UpdateOrgDomain(ctx context.Context, orgID, domainID, domain string) (*OrgDomain, error) {
	var out OrgDomain
	body := map[string]interface{}{"domain": domain}
	return &out, c.put(ctx, "/api/v1/orgs/"+orgID+"/auth/domains/"+domainID, body, &out)
}

func (c *Client) DeleteOrgDomain(ctx context.Context, orgID, domainID string) error {
	return c.del(ctx, "/api/v1/orgs/"+orgID+"/auth/domains/"+domainID)
}

func (c *Client) ListMappingOperators(ctx context.Context) ([]MappingOperator, error) {
	var out struct {
		Operators []MappingOperator `json:"operators"`
	}
	return out.Operators, c.get(ctx, "/api/v1/auth/mapping-operators", &out)
}
