// REST DTO types matching uigraph-api JSON responses.
package client

import (
	"encoding/json"
	"time"
)

// ── Auth ─────────────────────────────────────────────────────────────────────

type MeResponse struct {
	UserID       string `json:"userId"`
	OrgID        string `json:"orgId"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Login        string `json:"login"`
	Kind         string `json:"kind"`
	Role         string `json:"role"`
	AuthProvider string `json:"authProvider"`
}

type OrgSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

// ── Users (server admin) ────────────────────────────────────────────────────

type User struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Name       string     `json:"name"`
	Login      string     `json:"login"`
	Disabled   bool       `json:"disabled"`
	Role       string     `json:"role"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// ── SSO ───────────────────────────────────────────────────────────────────────

type OAuthProvider struct {
	ID             string    `json:"id"`
	ProviderName   string    `json:"providerName"`
	Type           string    `json:"type"`
	DisplayName    string    `json:"displayName"`
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

// ── Org ───────────────────────────────────────────────────────────────────────

type Org struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Member struct {
	UserID    string    `json:"userId"`
	OrgID     string    `json:"orgId"`
	Role      string    `json:"role"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Team struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"orgId"`
	Name       string    `json:"name"`
	Email      string    `json:"email,omitempty"`
	ExternalID string    `json:"externalId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TeamMember struct {
	TeamID     string    `json:"teamId"`
	UserID     string    `json:"userId"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Invitation struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgId"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Code      string     `json:"code"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type ServiceAccount struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"orgId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	Disabled    bool      `json:"disabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ServiceAccountToken struct {
	ID               string     `json:"id"`
	ServiceAccountID string     `json:"serviceAccountId"`
	Name             string     `json:"name"`
	Prefix           string     `json:"prefix"`
	ExpiresAt        *time.Time `json:"expiresAt,omitempty"`
	LastUsedAt       *time.Time `json:"lastUsedAt,omitempty"`
	Revoked          bool       `json:"revoked"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type CreatedToken struct {
	ServiceAccountToken
	Token string `json:"token"`
}

// ── Content ───────────────────────────────────────────────────────────────────

type Folder struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgId"`
	ParentID  *string    `json:"parentId,omitempty"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Order     float64    `json:"order"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type Diagram struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"orgId"`
	FolderID    *string    `json:"folderId,omitempty"`
	TeamID      *string    `json:"teamId,omitempty"`
	Name        string     `json:"name"`
	ContentKey  string     `json:"contentKey"`
	ContentHash string     `json:"contentHash"`
	Source      *string    `json:"source,omitempty"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type DiagramVersion struct {
	ID            string    `json:"id"`
	DiagramID     string    `json:"diagramId"`
	VersionNumber int       `json:"versionNumber"`
	Label         *string   `json:"label,omitempty"`
	ContentKey    string    `json:"contentKey"`
	ContentHash   string    `json:"contentHash"`
	IsAutoVersion bool      `json:"isAutoVersion"`
	Source        *string   `json:"source,omitempty"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type UIMap struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"orgId"`
	FolderID    *string    `json:"folderId,omitempty"`
	TeamID      *string    `json:"teamId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type Frame struct {
	ID                    string     `json:"id"`
	MapID                 string     `json:"mapId"`
	OrgID                 string     `json:"orgId"`
	ParentFrameID         *string    `json:"parentFrameId,omitempty"`
	Name                  string     `json:"name"`
	Description           string     `json:"description"`
	TemplateType          string     `json:"templateType"`
	ScreenshotKey         *string    `json:"screenshotKey,omitempty"`
	ScreenshotContentHash *string    `json:"screenshotContentHash,omitempty"`
	Status                string     `json:"status"`
	Order                 float64    `json:"order"`
	Source                *string    `json:"source,omitempty"`
	CreatedBy             string     `json:"createdBy"`
	UpdatedBy             *string    `json:"updatedBy,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty"`
}

type FocalPoint struct {
	ID         string     `json:"id"`
	FrameID    string     `json:"frameId"`
	OrgID      string     `json:"orgId"`
	Name       string     `json:"name"`
	LocationX  float64    `json:"locationX"`
	LocationY  float64    `json:"locationY"`
	Visibility string     `json:"visibility"`
	IsActive   bool       `json:"isActive"`
	CreatedBy  string     `json:"createdBy"`
	UpdatedBy  *string    `json:"updatedBy,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

type Canvas struct {
	MapID          string          `json:"mapId"`
	OrgID          string          `json:"orgId"`
	Zoom           float64         `json:"zoom"`
	NavigationX    float64         `json:"navigationX"`
	NavigationY    float64         `json:"navigationY"`
	FramePositions json.RawMessage `json:"framePositions"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// ── Catalog ───────────────────────────────────────────────────────────────────

type Service struct {
	ID              string          `json:"id"`
	OrgID           string          `json:"orgId"`
	FolderID        *string         `json:"folderId,omitempty"`
	TeamID          *string         `json:"teamId,omitempty"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	Description     string          `json:"description"`
	Status          string          `json:"status"`
	Tier            string          `json:"tier"`
	Category        string          `json:"category"`
	Language        string          `json:"language"`
	GitRepoURL      *string         `json:"gitRepoUrl,omitempty"`
	JiraProjectURL  *string         `json:"jiraProjectUrl,omitempty"`
	SlackChannelURL *string         `json:"slackChannelUrl,omitempty"`
	LastCommitSha   *string         `json:"lastCommitSha,omitempty"`
	Labels          []string        `json:"labels"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedBy       string          `json:"createdBy"`
	UpdatedBy       *string         `json:"updatedBy,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type APIGroup struct {
	ID        string     `json:"id"`
	ServiceID string     `json:"serviceId"`
	OrgID     string     `json:"orgId"`
	Name      string     `json:"name"`
	Version   string     `json:"version"`
	Label     *string    `json:"label,omitempty"`
	Protocol  string     `json:"protocol"`
	SpecKey   *string    `json:"specKey,omitempty"`
	SpecHash  *string    `json:"specHash,omitempty"`
	CreatedBy string     `json:"createdBy"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type APIGroupVersion struct {
	ID            string    `json:"id"`
	APIGroupID    string    `json:"apiGroupId"`
	VersionNumber int       `json:"versionNumber"`
	Label         *string   `json:"label,omitempty"`
	SpecKey       string    `json:"specKey"`
	SpecHash      string    `json:"specHash"`
	IsAutoVersion bool      `json:"isAutoVersion"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type APIEndpoint struct {
	ID          string          `json:"id"`
	APIGroupID  string          `json:"apiGroupId"`
	ServiceID   string          `json:"serviceId"`
	OrgID       string          `json:"orgId"`
	OperationID string          `json:"operationId"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	Summary     string          `json:"summary"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Parameters  json.RawMessage `json:"parameters"`
	RequestBody json.RawMessage `json:"requestBody"`
	Responses   json.RawMessage `json:"responses"`
	Order       float64         `json:"order"`
	CreatedBy   string          `json:"createdBy"`
	UpdatedBy   *string         `json:"updatedBy,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
