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
	TeamID    *string    `json:"teamId,omitempty"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Order     float64    `json:"order"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type Diagram struct {
	ID                 string     `json:"id"`
	OrgID              string     `json:"orgId"`
	FolderID           *string    `json:"folderId,omitempty"`
	TeamID             *string    `json:"teamId,omitempty"`
	Name               string     `json:"name"`
	ContentKey         string     `json:"contentKey"`
	ContentHash        string     `json:"contentHash"`
	PreviewAssetID     *string    `json:"previewAssetId,omitempty"`
	PreviewContentHash *string    `json:"previewContentHash,omitempty"`
	Source             *string    `json:"source,omitempty"`
	CreatedBy          string     `json:"createdBy"`
	UpdatedBy          *string    `json:"updatedBy,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt,omitempty"`
}

type FlowDiagramComponentField struct {
	FlowDiagramComponentFieldID string   `json:"flowDiagramComponentFieldId"`
	Label                       string   `json:"label"`
	Type                        string   `json:"type"`
	Required                    bool     `json:"required"`
	Readonly                    *bool    `json:"readonly,omitempty"`
	Options                     []string `json:"options,omitempty"`
	Order                       int      `json:"order"`
}

type FlowDiagramComponent struct {
	ComponentID                string                      `json:"componentId"`
	Type                       string                      `json:"type"`
	Name                       string                      `json:"name"`
	Description                string                      `json:"description"`
	Category                   string                      `json:"category"`
	Tags                       []string                    `json:"tags"`
	Slug                       string                      `json:"slug"`
	PreviewImageJpg            string                      `json:"previewImageJpg"`
	IsActive                   bool                        `json:"isActive"`
	Order                      int                         `json:"order"`
	OrganizationID             *string                     `json:"organizationId,omitempty"`
	FlowDiagramComponentFields []FlowDiagramComponentField `json:"flowDiagramComponentFields"`
}

type FlowComponents struct {
	Components       []FlowDiagramComponent `json:"components"`
	CustomComponents []FlowDiagramComponent `json:"customComponents"`
}

type ComponentField struct {
	ComponentFieldID string   `json:"componentFieldId"`
	Label            string   `json:"label"`
	Type             string   `json:"type"`
	Required         bool     `json:"required"`
	Readonly         *bool    `json:"readonly,omitempty"`
	Options          []string `json:"options,omitempty"`
	Order            int      `json:"order"`
}

type Component struct {
	ComponentID     string           `json:"componentId"`
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Category        string           `json:"category"`
	Tags            []string         `json:"tags"`
	Slug            string           `json:"slug"`
	PreviewImageJpg string           `json:"previewImageJpg"`
	IsActive        bool             `json:"isActive"`
	Order           int              `json:"order"`
	ComponentFields []ComponentField `json:"componentFields"`
}

type Components struct {
	Components       []Component `json:"components"`
	CustomComponents []Component `json:"customComponents"`
}

type DiagramImage struct {
	DiagramImageID string    `json:"diagramImageId"`
	DiagramID      string    `json:"diagramId"`
	OrgID          string    `json:"orgId"`
	AssetID        string    `json:"assetId"`
	FileName       *string   `json:"fileName,omitempty"`
	Order          int       `json:"order"`
	CreatedBy      string    `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
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
	ScreenshotAssetID     *string    `json:"screenshotAssetId,omitempty"`
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

type FrameGroup struct {
	ID          string     `json:"id"`
	FrameID     string     `json:"frameId"`
	OrgID       string     `json:"orgId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LocationX   float64    `json:"locationX"`
	LocationY   float64    `json:"locationY"`
	Width       float64    `json:"width"`
	Height      float64    `json:"height"`
	Order       float64    `json:"order"`
	IsActive    bool       `json:"isActive"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type FrameLink struct {
	ID            string     `json:"id"`
	FrameID       string     `json:"frameId"`
	OrgID         string     `json:"orgId"`
	Kind          string     `json:"kind"`
	TargetFrameID *string    `json:"targetFrameId,omitempty"`
	TargetMapID   *string    `json:"targetMapId,omitempty"`
	Label         string     `json:"label"`
	LocationX     float64    `json:"locationX"`
	LocationY     float64    `json:"locationY"`
	IsActive      bool       `json:"isActive"`
	CreatedBy     string     `json:"createdBy"`
	UpdatedBy     *string    `json:"updatedBy,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

type FocalPointMeta struct {
	ID                   string          `json:"id"`
	FocalPointID         string          `json:"focalPointId"`
	OrgID                string          `json:"orgId"`
	FrameID              string          `json:"frameId"`
	ComponentID          string          `json:"componentId"`
	ComponentLinkID      *string         `json:"componentLinkId,omitempty"`
	ComponentImages      json.RawMessage `json:"componentImages"`
	ComponentFlowDiagram *string         `json:"componentFlowDiagram,omitempty"`
	ComponentModalFields json.RawMessage `json:"componentModalFields"`
	CreatedBy            string          `json:"createdBy"`
	UpdatedBy            *string         `json:"updatedBy,omitempty"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
	DeletedAt            *time.Time      `json:"deletedAt,omitempty"`
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

type ServiceStats struct {
	ServiceID     string `json:"serviceId"`
	EndpointCount int    `json:"endpointCount"`
	DiagramCount  int    `json:"diagramCount"`
	DocCount      int    `json:"docCount"`
	DBTableCount  int    `json:"dbTableCount"`
	TestCaseCount int    `json:"testCaseCount"`
}

type APIGroup struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"serviceId"`
	OrgID     string    `json:"orgId"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Label     *string   `json:"label,omitempty"`
	Protocol  string    `json:"protocol"`
	SpecKey   *string   `json:"specKey,omitempty"`
	SpecHash  *string   `json:"specHash,omitempty"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

type ServiceDoc struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"serviceId"`
	OrgID       string    `json:"orgId"`
	FileKey     string    `json:"fileKey"`
	FileName    string    `json:"fileName"`
	FileType    string    `json:"fileType"`
	Description string    `json:"description"`
	ContentHash string    `json:"contentHash"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   *string   `json:"updatedBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ServiceDiagram struct {
	ServiceID string    `json:"serviceId"`
	DiagramID string    `json:"diagramId"`
	OrgID     string    `json:"orgId"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Diagram   *Diagram  `json:"diagram,omitempty"`
}

type ServiceDB struct {
	ID         string          `json:"id"`
	ServiceID  string          `json:"serviceId"`
	OrgID      string          `json:"orgId"`
	DBName     string          `json:"dbName"`
	DBType     string          `json:"dbType"`
	Dialect    string          `json:"dialect"`
	SchemaJSON json.RawMessage `json:"schemaJson"`
	Source     *string         `json:"source,omitempty"`
	SourceTS   *time.Time      `json:"sourceTs,omitempty"`
	CreatedBy  string          `json:"createdBy"`
	UpdatedBy  *string         `json:"updatedBy,omitempty"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type ServiceDBVersion struct {
	ID            string          `json:"id"`
	ServiceDBID   string          `json:"serviceDbId"`
	VersionNumber int             `json:"versionNumber"`
	Label         *string         `json:"label,omitempty"`
	SchemaJSON    json.RawMessage `json:"schemaJson"`
	Source        *string         `json:"source,omitempty"`
	SourceTS      *time.Time      `json:"sourceTs,omitempty"`
	IsAutoVersion bool            `json:"isAutoVersion"`
	CreatedBy     string          `json:"createdBy"`
	CreatedAt     time.Time       `json:"createdAt"`
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

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Assertion struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AuthConfig struct {
	Type          string  `json:"type"`
	BearerToken   *string `json:"bearerToken,omitempty"`
	APIKeyHeader  *string `json:"apiKeyHeader,omitempty"`
	APIKeyValue   *string `json:"apiKeyValue,omitempty"`
	BasicUsername *string `json:"basicUsername,omitempty"`
	BasicPassword *string `json:"basicPassword,omitempty"`
}

type TestCaseStep struct {
	Order          int    `json:"order"`
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

type ManualTestCase struct {
	Preconditions   *string        `json:"preconditions,omitempty"`
	TestData        *string        `json:"testData,omitempty"`
	Steps           []TestCaseStep `json:"steps,omitempty"`
	ExpectedOutcome *string        `json:"expectedOutcome,omitempty"`
	Postconditions  *string        `json:"postconditions,omitempty"`
}

type APITestCase struct {
	HTTPMethod         string      `json:"httpMethod"`
	APISpecID          *string     `json:"apiSpecId,omitempty"`
	OperationID        *string     `json:"operationId,omitempty"`
	Auth               *AuthConfig `json:"auth,omitempty"`
	RequestHeaders     []KeyValue  `json:"requestHeaders,omitempty"`
	QueryParams        []KeyValue  `json:"queryParams,omitempty"`
	RequestBody        *string     `json:"requestBody,omitempty"`
	ExpectedStatusCode *int        `json:"expectedStatusCode,omitempty"`
	MaxResponseTimeMs  *int        `json:"maxResponseTimeMs,omitempty"`
	ResponseBody       *string     `json:"responseBody,omitempty"`
	Assertions         []Assertion `json:"assertions,omitempty"`
}

type GraphQLTestCase struct {
	OperationType string      `json:"operationType"`
	OperationName *string     `json:"operationName,omitempty"`
	Query         string      `json:"query"`
	Variables     *string     `json:"variables,omitempty"`
	ResponseBody  *string     `json:"responseBody,omitempty"`
	Assertions    []Assertion `json:"assertions,omitempty"`
	ExpectError   bool        `json:"expectError"`
}

type DatabaseTestCase struct {
	Dialect       string      `json:"dialect"`
	SchemaID      *string     `json:"schemaId,omitempty"`
	Query         string      `json:"query"`
	Assertions    []Assertion `json:"assertions,omitempty"`
	SetupQuery    *string     `json:"setupQuery,omitempty"`
	TeardownQuery *string     `json:"teardownQuery,omitempty"`
}

type GRPCTestCase struct {
	ServiceName    string      `json:"serviceName"`
	MethodName     string      `json:"methodName"`
	CallMode       string      `json:"callMode"`
	ProtoFileID    *string     `json:"protoFileId,omitempty"`
	ServerAddress  *string     `json:"serverAddress,omitempty"`
	RequestMessage *string     `json:"requestMessage,omitempty"`
	Metadata       []KeyValue  `json:"metadata,omitempty"`
	ExpectedStatus string      `json:"expectedStatus"`
	DeadlineMs     *int        `json:"deadlineMs,omitempty"`
	ResponseBody   *string     `json:"responseBody,omitempty"`
	Assertions     []Assertion `json:"assertions,omitempty"`
	UseTLS         bool        `json:"useTLS"`
	ExpectError    bool        `json:"expectError"`
}

type TestPack struct {
	TestPackID string     `json:"testPackId"`
	ServiceID  string     `json:"serviceId"`
	OrgID      string     `json:"orgId"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	CreatedBy  string     `json:"createdBy"`
	UpdatedBy  *string    `json:"updatedBy,omitempty"`
	DeletedBy  *string    `json:"deletedBy,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

type TestCase struct {
	TestCaseID            string            `json:"testCaseId"`
	TestPackID            string            `json:"testPackId"`
	ServiceID             string            `json:"serviceId"`
	OrgID                 string            `json:"orgId"`
	Title                 string            `json:"title"`
	Order                 float64           `json:"order"`
	Type                  string            `json:"type"`
	Description           *string           `json:"description,omitempty"`
	Priority              *string           `json:"priority,omitempty"`
	Labels                []string          `json:"labels,omitempty"`
	LinkedTicket          *string           `json:"linkedTicket,omitempty"`
	EstimatedDurationMins *int              `json:"estimatedDurationMins,omitempty"`
	TestOwner             *string           `json:"testOwner,omitempty"`
	LinkedMapNodeID       *string           `json:"linkedMapNodeId,omitempty"`
	IsCritical            bool              `json:"isCritical"`
	EvidenceRequired      bool              `json:"evidenceRequired"`
	Manual                *ManualTestCase   `json:"manual,omitempty"`
	API                   *APITestCase      `json:"api,omitempty"`
	GraphQL               *GraphQLTestCase  `json:"graphql,omitempty"`
	Database              *DatabaseTestCase `json:"database,omitempty"`
	GRPC                  *GRPCTestCase     `json:"grpc,omitempty"`
	Status                string            `json:"status"`
	Version               int               `json:"version"`
	BaselineRunResultID   *string           `json:"baselineRunResultId,omitempty"`
	Dependencies          []string          `json:"dependencies,omitempty"`
	CreatedBy             string            `json:"createdBy"`
	UpdatedBy             *string           `json:"updatedBy,omitempty"`
	DeletedBy             *string           `json:"deletedBy,omitempty"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
	DeletedAt             *time.Time        `json:"deletedAt,omitempty"`
}

type TestRun struct {
	TestRunID     string     `json:"testRunId"`
	TestPackID    string     `json:"testPackId"`
	ServiceID     string     `json:"serviceId"`
	OrgID         string     `json:"orgId"`
	Environment   string     `json:"environment"`
	ReleaseLabel  *string    `json:"releaseLabel,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	StartedBy     *string    `json:"startedBy,omitempty"`
	ExecutedBy    string     `json:"executedBy"`
	ExecutedAt    time.Time  `json:"executedAt"`
	OverallStatus string     `json:"overallStatus"`
}

type TestRunSummary struct {
	TestRunID     string     `json:"testRunId"`
	TestPackID    string     `json:"testPackId"`
	ServiceID     string     `json:"serviceId"`
	Environment   string     `json:"environment"`
	ReleaseLabel  *string    `json:"releaseLabel,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	StartedBy     *string    `json:"startedBy,omitempty"`
	ExecutedBy    string     `json:"executedBy"`
	ExecutedAt    time.Time  `json:"executedAt"`
	OverallStatus string     `json:"overallStatus"`
	PassedCount   int        `json:"passedCount"`
	FailedCount   int        `json:"failedCount"`
	SkippedCount  int        `json:"skippedCount"`
	BlockedCount  int        `json:"blockedCount"`
}

type TestRunResult struct {
	TestRunResultID string    `json:"testRunResultId"`
	TestRunID       string    `json:"testRunId"`
	TestCaseID      string    `json:"testCaseId"`
	ServiceID       string    `json:"serviceId"`
	OrgID           string    `json:"orgId"`
	Status          string    `json:"status"`
	BlockedReason   *string   `json:"blockedReason,omitempty"`
	ResponseStatus  *int      `json:"responseStatus,omitempty"`
	ResponseBody    *string   `json:"responseBody,omitempty"`
	ResponseTimeMs  *int64    `json:"responseTimeMs,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	ScreenshotURLs  []string  `json:"screenshotUrls,omitempty"`
	ExecutedAt      time.Time `json:"executedAt"`
	ExecutedBy      string    `json:"executedBy"`
}
