package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

type authClient interface {
	Me(ctx context.Context) (*uigraphapi.MeResponse, error)
	MyOrgs(ctx context.Context) ([]uigraphapi.OrgSummary, error)
	SwitchOrg(ctx context.Context, orgID string) error
}

type orgClient interface {
	ListOrgs(ctx context.Context) ([]uigraphapi.Org, error)
	GetOrg(ctx context.Context, id string) (*uigraphapi.Org, error)
	CreateOrg(ctx context.Context, body map[string]interface{}) (*uigraphapi.Org, error)
	UpdateOrg(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.Org, error)
	DeleteOrg(ctx context.Context, id string) error
	ListMembers(ctx context.Context, orgID string) ([]uigraphapi.Member, error)
	AddMember(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Member, error)
	UpdateMember(ctx context.Context, orgID, userID string, body map[string]interface{}) (*uigraphapi.Member, error)
	RemoveMember(ctx context.Context, orgID, userID string) error
	ListTeams(ctx context.Context, orgID string) ([]uigraphapi.Team, error)
	GetTeam(ctx context.Context, orgID, teamID string) (*uigraphapi.Team, error)
	CreateTeam(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Team, error)
	UpdateTeam(ctx context.Context, orgID, teamID string, body map[string]interface{}) (*uigraphapi.Team, error)
	DeleteTeam(ctx context.Context, orgID, teamID string) error
	ListTeamMembers(ctx context.Context, orgID, teamID string) ([]uigraphapi.TeamMember, error)
	AddTeamMember(ctx context.Context, orgID, teamID string, body map[string]interface{}) error
	RemoveTeamMember(ctx context.Context, orgID, teamID, userID string) error
	ListServiceAccounts(ctx context.Context, orgID string) ([]uigraphapi.ServiceAccount, error)
	GetServiceAccount(ctx context.Context, orgID, id string) (*uigraphapi.ServiceAccount, error)
	CreateServiceAccount(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.ServiceAccount, error)
	UpdateServiceAccount(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, orgID, id string) error
	ListServiceAccountTokens(ctx context.Context, orgID, saID string) ([]uigraphapi.ServiceAccountToken, error)
	CreateServiceAccountToken(ctx context.Context, orgID, saID string, body map[string]interface{}) (*uigraphapi.CreatedToken, error)
	RevokeServiceAccountToken(ctx context.Context, orgID, saID, tokenID string) error
	ListServiceAccountScopes(ctx context.Context, orgID string) ([]string, error)
}

type adminClient interface {
	GetServerOverview(ctx context.Context) (*uigraphapi.ServerOverview, error)
	GetServerConfig(ctx context.Context) (*uigraphapi.ServerConfig, error)
	ServerListOrgs(ctx context.Context) ([]uigraphapi.Org, error)
	ServerCreateOrg(ctx context.Context, body map[string]interface{}) (*uigraphapi.Org, error)
	ServerUpdateOrg(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.Org, error)
	ServerDeleteOrg(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]uigraphapi.User, error)
	GetUser(ctx context.Context, id string) (*uigraphapi.User, error)
	CreateUser(ctx context.Context, body map[string]interface{}) (*uigraphapi.User, error)
	UpdateUser(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.User, error)
	DisableUser(ctx context.Context, id string) error
	ListOAuthProviders(ctx context.Context) ([]uigraphapi.OAuthProvider, error)
	UpsertOAuthProvider(ctx context.Context, provider string, body map[string]interface{}) error
	DeleteOAuthProvider(ctx context.Context, provider string) error
	ListRoleMappings(ctx context.Context) ([]uigraphapi.RoleMapping, error)
	CreateRoleMapping(ctx context.Context, body map[string]interface{}) error
	DeleteRoleMapping(ctx context.Context, id string) error
	GetLDAP(ctx context.Context) (*uigraphapi.LDAPConfig, error)
	UpsertLDAP(ctx context.Context, body map[string]interface{}) error
	DeleteLDAP(ctx context.Context) error
	GetSAML(ctx context.Context) (*uigraphapi.SAMLConfig, error)
	UpsertSAML(ctx context.Context, body map[string]interface{}) error
	GetSCIM(ctx context.Context) (*uigraphapi.SCIMConfig, error)
}

type folderClient interface {
	ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]uigraphapi.Folder, error)
	GetFolder(ctx context.Context, orgID, id string) (*uigraphapi.Folder, error)
	CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Folder, error)
	UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Folder, error)
	DeleteFolder(ctx context.Context, orgID, id string) error
}

type diagramClient interface {
	ListDiagrams(ctx context.Context, orgID, folderID string) ([]uigraphapi.Diagram, error)
	GetDiagram(ctx context.Context, orgID, id string) (*uigraphapi.Diagram, error)
	GetDiagramContent(ctx context.Context, orgID, id string) (string, error)
	CreateDiagram(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Diagram, error)
	UpdateDiagram(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Diagram, error)
	DeleteDiagram(ctx context.Context, orgID, id string) error
	ListDiagramImages(ctx context.Context, orgID, diagramID string) ([]uigraphapi.DiagramImage, error)
	CreateDiagramImage(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (*uigraphapi.DiagramImage, error)
	SyncDiagram(ctx context.Context, orgID string, body map[string]interface{}) (map[string]interface{}, error)
	ListDiagramVersions(ctx context.Context, orgID, diagramID string) ([]uigraphapi.DiagramVersion, error)
	CreateDiagramVersion(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (*uigraphapi.DiagramVersion, error)
	GetDiagramVersionContent(ctx context.Context, orgID, diagramID, versionID string) (string, error)
	RestoreDiagramVersion(ctx context.Context, orgID, diagramID, versionID string) (*uigraphapi.Diagram, error)
	PrepareDiagramThumbnailUpload(ctx context.Context, orgID, diagramID string) (*uigraphapi.DiagramThumbnailUpload, error)
	ConfirmDiagramThumbnailUpload(ctx context.Context, orgID, diagramID, contentHash string) error
}

type docsClient interface {
	ListDocs(ctx context.Context, orgID, folderID string) ([]uigraphapi.Doc, error)
	GetDoc(ctx context.Context, orgID, id string) (*uigraphapi.Doc, error)
	CreateDoc(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Doc, error)
	UpdateDoc(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Doc, error)
	DeleteDoc(ctx context.Context, orgID, id string) error
}

type componentClient interface {
	ListFlowDiagramComponents(ctx context.Context, orgID string) (*uigraphapi.FlowComponents, error)
	ListComponents(ctx context.Context, orgID string) (*uigraphapi.Components, error)
	CreateCustomComponent(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Component, error)
	UpdateCustomComponent(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Component, error)
	DeleteCustomComponent(ctx context.Context, orgID, id string) error
}

type uimapClient interface {
	ListMaps(ctx context.Context, orgID, folderID string) ([]uigraphapi.UIMap, error)
	GetMap(ctx context.Context, orgID, id string) (*uigraphapi.UIMap, error)
	CreateMap(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.UIMap, error)
	UpdateMap(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.UIMap, error)
	DeleteMap(ctx context.Context, orgID, id string) error
	ListFrames(ctx context.Context, orgID, mapID string) ([]uigraphapi.Frame, error)
	GetFrame(ctx context.Context, orgID, mapID, id string) (*uigraphapi.Frame, error)
	GetFrameByID(ctx context.Context, orgID, id string) (*uigraphapi.Frame, error)
	CreateFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*uigraphapi.Frame, error)
	UpdateFrame(ctx context.Context, orgID, mapID, id string, body map[string]interface{}) (*uigraphapi.Frame, error)
	DeleteFrame(ctx context.Context, orgID, mapID, id string) error
	SyncFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (map[string]interface{}, error)
	ListFocalPoints(ctx context.Context, orgID, mapID, frameID string) ([]uigraphapi.FocalPoint, error)
	GetFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) (*uigraphapi.FocalPoint, error)
	CreateFocalPoint(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*uigraphapi.FocalPoint, error)
	UpdateFocalPoint(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*uigraphapi.FocalPoint, error)
	DeleteFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) error
	GetCanvas(ctx context.Context, orgID, mapID string) (*uigraphapi.Canvas, error)
	UpsertCanvas(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*uigraphapi.Canvas, error)
	ListFrameGroups(ctx context.Context, orgID, mapID, frameID string) ([]uigraphapi.FrameGroup, error)
	CreateFrameGroup(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*uigraphapi.FrameGroup, error)
	UpdateFrameGroup(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*uigraphapi.FrameGroup, error)
	DeleteFrameGroup(ctx context.Context, orgID, mapID, frameID, id string) error
	ListFrameLinks(ctx context.Context, orgID, mapID, frameID string) ([]uigraphapi.FrameLink, error)
	CreateFrameLink(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*uigraphapi.FrameLink, error)
	UpdateFrameLink(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*uigraphapi.FrameLink, error)
	DeleteFrameLink(ctx context.Context, orgID, mapID, frameID, id string) error
	ListFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string) ([]uigraphapi.FocalPointMeta, error)
	ListFocalPointMetaByLink(ctx context.Context, orgID, linkKey, linkValue string) ([]uigraphapi.FocalPointMeta, error)
	CreateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string, body map[string]interface{}) (*uigraphapi.FocalPointMeta, error)
	UpdateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string, body map[string]interface{}) (*uigraphapi.FocalPointMeta, error)
	DeleteFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string) error
}

type catalogClient interface {
	ListServices(ctx context.Context, orgID, folderID, teamID string) ([]uigraphapi.Service, error)
	GetService(ctx context.Context, orgID, id string) (*uigraphapi.Service, error)
	CreateService(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Service, error)
	UpdateService(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Service, error)
	DeleteService(ctx context.Context, orgID, id string) error
	ListServiceStats(ctx context.Context, orgID string, serviceID *string) ([]uigraphapi.ServiceStats, error)
	ListAPIGroups(ctx context.Context, orgID, serviceID string) ([]uigraphapi.APIGroup, error)
	GetAPIGroup(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.APIGroup, error)
	GetAPIGroupSpec(ctx context.Context, orgID, serviceID, apiGroupID, versionID string) (*uigraphapi.APIGroupSpec, error)
	CreateAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.APIGroup, error)
	UpdateAPIGroup(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.APIGroup, error)
	DeleteAPIGroup(ctx context.Context, orgID, serviceID, id string) error
	SyncAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (map[string]interface{}, error)
	ListAPIGroupVersions(ctx context.Context, orgID, serviceID, apiGroupID string) ([]uigraphapi.APIGroupVersion, error)
	ListServiceDocs(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceDoc, error)
	CreateServiceDoc(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceDoc, error)
	DeleteServiceDoc(ctx context.Context, orgID, serviceID, docID string) error
	ListServiceDiagrams(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceDiagram, error)
	CreateServiceDiagram(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceDiagram, error)
	DeleteServiceDiagram(ctx context.Context, orgID, serviceID, diagramID string) error
	ListServiceDBs(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceDB, error)
	GetServiceDB(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.ServiceDB, error)
	CreateServiceDB(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceDB, error)
	UpdateServiceDB(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.ServiceDB, error)
	DeleteServiceDB(ctx context.Context, orgID, serviceID, id string) error
	ListServiceDBVersions(ctx context.Context, orgID, serviceID, serviceDBID string) ([]uigraphapi.ServiceDBVersion, error)
	CreateServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID string, body map[string]interface{}) (*uigraphapi.ServiceDBVersion, error)
	RestoreServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID, versionID string) (*uigraphapi.ServiceDB, error)
	ListAPIEndpoints(ctx context.Context, orgID, serviceID, apiGroupID string) ([]uigraphapi.APIEndpoint, error)
	GetAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) (*uigraphapi.APIEndpoint, error)
	CreateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID string, body map[string]interface{}) (*uigraphapi.APIEndpoint, error)
	UpdateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string, body map[string]interface{}) (*uigraphapi.APIEndpoint, error)
	DeleteAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) error
}

type testPackClient interface {
	ListTestPacks(ctx context.Context, orgID, serviceID string) ([]uigraphapi.TestPack, error)
	CreateTestPack(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestPack, error)
	UpdateTestPack(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestPack, error)
	DeleteTestPack(ctx context.Context, orgID, serviceID, id string) error
	ListTestCases(ctx context.Context, orgID, serviceID string, testPackID *string) ([]uigraphapi.TestCase, error)
	CreateTestCase(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestCase, error)
	UpdateTestCase(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestCase, error)
	DeleteTestCase(ctx context.Context, orgID, serviceID, id string) error
	GetTestRun(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.TestRun, error)
	ListTestRuns(ctx context.Context, orgID, serviceID string, testPackID *string) ([]uigraphapi.TestRun, error)
	ListTestRunsSummary(ctx context.Context, orgID, serviceID string, testPackID, environment, status, executedBy *string, fromDate, toDate *time.Time) ([]uigraphapi.TestRunSummary, error)
	CreateTestRun(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestRun, error)
	UpdateTestRun(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestRun, error)
	ListTestRunResults(ctx context.Context, orgID, serviceID, testRunID string) ([]uigraphapi.TestRunResult, error)
	CreateTestRunResult(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestRunResult, error)
	UpdateTestRunResult(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestRunResult, error)
}

type actorClient interface {
	ResolveActors(ctx context.Context, orgID string, ids []string) (map[string]*uigraphapi.Actor, error)
	ResolveAssetURLs(ctx context.Context, orgID string, ids []string) (map[string]string, error)
	CreateAssetUpload(ctx context.Context, orgID string) (*uigraphapi.AssetUpload, error)
}

type commentClient interface {
	ListComments(ctx context.Context, orgID, resourceID string) ([]uigraphapi.Comment, error)
	CreateComment(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Comment, error)
	UpdateComment(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Comment, error)
	DeleteComment(ctx context.Context, orgID, id string) error
}

// Resolver is the root dependency-injection struct for all resolvers. Each
// field is the minimal interface its domain's resolvers need — not the full
// *uigraphapi.Client — so tests can inject a narrow fake instead of mocking
// every REST method.
type Resolver struct {
	Auth       authClient
	OrgAPI     orgClient
	Admin      adminClient
	FolderAPI  folderClient
	DiagramAPI diagramClient
	DocAPI     docsClient
	Component  componentClient
	UIMapAPI   uimapClient
	Catalog    catalogClient
	TestPack   testPackClient
	Actor      actorClient
	CommentAPI commentClient
}
