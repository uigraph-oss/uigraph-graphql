package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (r *queryResolver) resolveSessionActors(ctx context.Context, orgID string, sessions []uigraphapi.AgentSession) (map[string]*uigraphapi.Actor, error) {
	ids := make([]string, 0, len(sessions))
	for _, s := range sessions {
		if s.UserID != nil {
			ids = append(ids, *s.UserID)
		} else if s.ServiceAccountID != nil {
			ids = append(ids, *s.ServiceAccountID)
		}
	}
	if len(ids) == 0 {
		return map[string]*uigraphapi.Actor{}, nil
	}
	return r.Resolver.Actor.ResolveActors(ctx, orgID, ids)
}

type authClient interface {
	Me(ctx context.Context) (*uigraphapi.MeResponse, error)
	MyOrgs(ctx context.Context) ([]uigraphapi.OrgSummary, error)
	SwitchOrg(ctx context.Context, orgID string) error
	PrepareUserAvatarUpload(ctx context.Context) (*uigraphapi.AssetUpload, error)
	SetMyAvatar(ctx context.Context) error
}

type orgClient interface {
	ListOrgs(ctx context.Context) ([]uigraphapi.Org, error)
	GetOrg(ctx context.Context, id string) (*uigraphapi.Org, error)
	CreateOrg(ctx context.Context, body map[string]interface{}) (*uigraphapi.Org, error)
	UpdateOrg(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.Org, error)
	DeleteOrg(ctx context.Context, id string) error
	CompleteOnboarding(ctx context.Context, orgID string) error
	PrepareOrgLogoUpload(ctx context.Context, orgID string) (*uigraphapi.AssetUpload, error)
	SetOrgLogo(ctx context.Context, orgID string) error
	RemoveOrgLogo(ctx context.Context, orgID string) error
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
	PrepareServiceAccountAvatarUpload(ctx context.Context, orgID, saID string) (*uigraphapi.AssetUpload, error)
	SetServiceAccountAvatar(ctx context.Context, orgID, saID string) error
}

type adminClient interface {
	GetServerOverview(ctx context.Context) (*uigraphapi.ServerOverview, error)
	GetServerConfig(ctx context.Context) (*uigraphapi.ServerConfig, error)
	ServerListOrgs(ctx context.Context) ([]uigraphapi.Org, error)
	ServerCreateOrg(ctx context.Context, body map[string]interface{}) (*uigraphapi.Org, error)
	ServerUpdateOrg(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.Org, error)
	ServerDeleteOrg(ctx context.Context, id string) error
	PrepareServerOrgLogoUpload(ctx context.Context, orgID string) (*uigraphapi.AssetUpload, error)
	SetServerOrgLogo(ctx context.Context, orgID string) error
	RemoveServerOrgLogo(ctx context.Context, orgID string) error
	ListUsers(ctx context.Context) ([]uigraphapi.User, error)
	GetUser(ctx context.Context, id string) (*uigraphapi.User, error)
	CreateUser(ctx context.Context, body map[string]interface{}) (*uigraphapi.User, error)
	UpdateUser(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.User, error)
	DisableUser(ctx context.Context, id string) error
	GetLDAP(ctx context.Context) (*uigraphapi.LDAPConfig, error)
	UpsertLDAP(ctx context.Context, body map[string]interface{}) error
	DeleteLDAP(ctx context.Context) error
	GetSCIM(ctx context.Context) (*uigraphapi.SCIMConfig, error)
}

type authProviderClient interface {
	ListAuthProviders(ctx context.Context, orgID string) ([]uigraphapi.AuthProvider, error)
	GetAuthProvider(ctx context.Context, orgID, slug string) (*uigraphapi.AuthProvider, error)
	CreateAuthProvider(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.AuthProvider, error)
	UpdateAuthProvider(ctx context.Context, orgID, slug string, body map[string]interface{}) (*uigraphapi.AuthProvider, error)
	DeleteAuthProvider(ctx context.Context, orgID, slug string) error
	PrepareAuthProviderIconUpload(ctx context.Context, orgID, slug string) (*uigraphapi.AssetUpload, error)
	SetAuthProviderIcon(ctx context.Context, orgID, slug string) error
	RemoveAuthProviderIcon(ctx context.Context, orgID, slug string) error
	GetAuthProviderSAMLMetadata(ctx context.Context, orgID, slug string) (*uigraphapi.SamlSpMetadata, error)
	ListAuthRoleMappings(ctx context.Context, orgID, providerSlug string) ([]uigraphapi.AuthRoleMapping, error)
	CreateAuthRoleMapping(ctx context.Context, orgID, providerSlug string, body map[string]interface{}) (*uigraphapi.AuthRoleMapping, error)
	UpdateAuthRoleMapping(ctx context.Context, orgID, providerSlug, mappingID string, body map[string]interface{}) (*uigraphapi.AuthRoleMapping, error)
	DeleteAuthRoleMapping(ctx context.Context, orgID, providerSlug, mappingID string) error
	ListOrgDomains(ctx context.Context, orgID string) ([]uigraphapi.OrgDomain, error)
	CreateOrgDomain(ctx context.Context, orgID, domain string) (*uigraphapi.OrgDomain, error)
	UpdateOrgDomain(ctx context.Context, orgID, domainID, domain string) (*uigraphapi.OrgDomain, error)
	DeleteOrgDomain(ctx context.Context, orgID, domainID string) error
	ListMappingOperators(ctx context.Context) ([]uigraphapi.MappingOperator, error)
}

type folderClient interface {
	ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]uigraphapi.Folder, error)
	GetFolder(ctx context.Context, orgID, id string) (*uigraphapi.Folder, error)
	CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Folder, error)
	UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Folder, error)
	DeleteFolder(ctx context.Context, orgID, id string) error
}

type diagramClient interface {
	ListDiagrams(ctx context.Context, orgID string, p uigraphapi.ListParams) ([]uigraphapi.Diagram, int, error)
	GetDiagram(ctx context.Context, orgID, id string, includeContent bool) (*uigraphapi.Diagram, error)
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
	GenerateDiagramThumbnail(ctx context.Context, orgID, diagramID string) error
}

type docsClient interface {
	ListDocs(ctx context.Context, orgID string, p uigraphapi.ListParams) ([]uigraphapi.Doc, int, error)
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
	ListMaps(ctx context.Context, orgID string, p uigraphapi.ListParams) ([]uigraphapi.UIMap, int, error)
	GetMap(ctx context.Context, orgID, id string) (*uigraphapi.UIMap, error)
	CreateMap(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.UIMap, error)
	UpdateMap(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.UIMap, error)
	DeleteMap(ctx context.Context, orgID, id string) error
	ListFrames(ctx context.Context, orgID, mapID string, p uigraphapi.ListParams) ([]uigraphapi.Frame, int, error)
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
	ListFocalPointMetaByLink(ctx context.Context, orgID, linkID string) ([]uigraphapi.FocalPointMeta, error)
	ListComponentLinkUsages(ctx context.Context, orgID, linkID string) ([]uigraphapi.ComponentLinkUsage, error)
	CreateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string, body map[string]interface{}) (*uigraphapi.FocalPointMeta, error)
	UpdateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string, body map[string]interface{}) (*uigraphapi.FocalPointMeta, error)
	DeleteFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string) error
}

type catalogClient interface {
	ListServices(ctx context.Context, orgID string, p uigraphapi.ListParams) ([]uigraphapi.Service, int, error)
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
	RestoreAPIGroupVersion(ctx context.Context, orgID, serviceID, apiGroupID, versionID string) (*uigraphapi.APIGroup, error)
	ListSavedQueryFolders(ctx context.Context, orgID, serviceID, dbID, scope string) ([]uigraphapi.SavedQueryFolder, error)
	CreateSavedQueryFolder(ctx context.Context, orgID, serviceID, dbID string, body map[string]interface{}) (*uigraphapi.SavedQueryFolder, error)
	DeleteSavedQueryFolder(ctx context.Context, orgID, serviceID, dbID, id string) error
	ListSavedQueries(ctx context.Context, orgID, serviceID, dbID, scope string) ([]uigraphapi.SavedQuery, error)
	CreateSavedQuery(ctx context.Context, orgID, serviceID, dbID string, body map[string]interface{}) (*uigraphapi.SavedQuery, error)
	UpdateSavedQuery(ctx context.Context, orgID, serviceID, dbID, id string, body map[string]interface{}) (*uigraphapi.SavedQuery, error)
	DeleteSavedQuery(ctx context.Context, orgID, serviceID, dbID, id string) error
	ListAPIEndpoints(ctx context.Context, orgID, serviceID, apiGroupID, versionID string) ([]uigraphapi.APIEndpoint, error)
	GetAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) (*uigraphapi.APIEndpoint, error)
	GetAPIEndpointByID(ctx context.Context, orgID, id string) (*uigraphapi.APIEndpoint, error)
	GetServiceDocByID(ctx context.Context, orgID, id string) (*uigraphapi.ServiceDoc, error)
	CreateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID string, body map[string]interface{}) (*uigraphapi.APIEndpoint, error)
	UpdateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string, body map[string]interface{}) (*uigraphapi.APIEndpoint, error)
	DeleteAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) error
}

type dependencyClient interface {
	ListDependencies(ctx context.Context, orgID, serviceID string, direction, criticality *string) ([]uigraphapi.Dependency, error)
	GetServiceDependencyGraph(ctx context.Context, orgID, serviceID string) ([]uigraphapi.Dependency, error)
	GetDependencyGraph(ctx context.Context, orgID string) ([]uigraphapi.Dependency, error)
	GetServiceImpact(ctx context.Context, orgID, serviceID string, direction *string, maxDepth *int) ([]uigraphapi.Dependency, error)
	UpdateServiceDependencies(ctx context.Context, orgID, serviceID string, body map[string]interface{}) ([]uigraphapi.Dependency, error)
}

type testPackClient interface {
	ListTestPacks(ctx context.Context, orgID, serviceID string) ([]uigraphapi.TestPack, error)
	GetTestPackByID(ctx context.Context, orgID, id string) (*uigraphapi.TestPack, error)
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

type chatClient interface {
	ListChatSessions(ctx context.Context, orgID string) ([]uigraphapi.ChatSession, error)
	CreateChatSession(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.ChatSession, error)
	GetChatSession(ctx context.Context, orgID, id string) (*uigraphapi.ChatSession, []uigraphapi.ChatMessage, error)
	UpdateChatSession(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.ChatSession, error)
	DeleteChatSession(ctx context.Context, orgID, id string) error
	ListChatMessages(ctx context.Context, orgID, sessionID string) ([]uigraphapi.ChatMessage, error)
	CreateChatMessage(ctx context.Context, orgID, sessionID string, body map[string]interface{}) (*uigraphapi.ChatMessage, error)
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

type costSavingsClient interface {
	GetSavingsSummary(ctx context.Context, orgID string, period, modelID *string) (*uigraphapi.SavingsSummary, error)
	GetSavingsTimeseries(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.DailySavings, error)
	GetSavingsByTool(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.ToolSavings, error)
	GetSavingsByClient(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.ClientSavings, error)
	GetSavingsByModel(ctx context.Context, orgID string, period *string) ([]uigraphapi.ModelSavings, error)
	GetSavingsByUser(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.UserSavings, error)
}

type agentSessionClient interface {
	GetAgentSessions(ctx context.Context, orgID string, sessionType, status, period *string, limit, offset *int) (*uigraphapi.AgentSessionPage, error)
	GetAgentSession(ctx context.Context, orgID, id string) (*uigraphapi.AgentSessionDetail, error)
	GetAgentSessionSummary(ctx context.Context, orgID string, period, sessionType *string) (*uigraphapi.AgentSessionSummary, error)
}

type mlStudioClient interface {
	ListMLProjects(ctx context.Context, orgID string) ([]uigraphapi.MLProject, error)
	GetMLProject(ctx context.Context, orgID, id string) (*uigraphapi.MLProject, error)
	CreateMLProject(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.MLProject, error)
	UpdateMLProject(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLProject, error)
	DeleteMLProject(ctx context.Context, orgID, id string) error
	ListMLModels(ctx context.Context, orgID, projectID string) ([]uigraphapi.MLModel, error)
	GetMLModel(ctx context.Context, orgID, id string) (*uigraphapi.MLModel, error)
	CreateMLModel(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.MLModel, error)
	UpdateMLModel(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLModel, error)
	UpdateMLModelInfo(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLModel, error)
	DeleteMLModel(ctx context.Context, orgID, id string) error
	ListMLModelVersions(ctx context.Context, orgID, modelID, projectID string) ([]uigraphapi.MLModelVersion, error)
	ListMLModelVersionsExplore(ctx context.Context, orgID string, query uigraphapi.MLModelVersionQuery) ([]uigraphapi.MLModelVersionExploreItem, int, error)
	GetMLModelVersion(ctx context.Context, orgID, id string) (*uigraphapi.MLModelVersion, error)
	ListMLExperiments(ctx context.Context, orgID, projectID string) ([]uigraphapi.MLExperiment, error)
	GetMLExperiment(ctx context.Context, orgID, id string) (*uigraphapi.MLExperiment, error)
	CreateMLExperiment(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.MLExperiment, error)
	UpdateMLExperiment(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLExperiment, error)
	DeleteMLExperiment(ctx context.Context, orgID, id string) error
	ListMLRuns(ctx context.Context, orgID string, query uigraphapi.MLRunQuery) ([]uigraphapi.MLRun, int, error)
	GetMLRun(ctx context.Context, orgID, id string) (*uigraphapi.MLRun, error)
	CreateMLRun(ctx context.Context, orgID, experimentID string, body map[string]interface{}) (*uigraphapi.MLRun, error)
	UpdateMLRun(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLRun, error)
	DeleteMLRun(ctx context.Context, orgID, id string) error
	ListMLArtifacts(ctx context.Context, orgID, runID string) ([]uigraphapi.MLArtifact, error)
	CreateMLArtifact(ctx context.Context, orgID, runID string, body map[string]interface{}) (*uigraphapi.MLArtifact, error)
	UpdateMLArtifact(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLArtifact, error)
	DeleteMLArtifact(ctx context.Context, orgID, id string) error
	ListMLDatasets(ctx context.Context, orgID, experimentID string) ([]uigraphapi.MLDataset, error)
	GetMLDataset(ctx context.Context, orgID, id string) (*uigraphapi.MLDataset, error)
	CreateMLDataset(ctx context.Context, orgID, experimentID string, body map[string]interface{}) (*uigraphapi.MLDataset, error)
	UpdateMLDataset(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLDataset, error)
	DeleteMLDataset(ctx context.Context, orgID, id string) error
	ListMLDeployments(ctx context.Context, orgID, modelID, versionID string) ([]uigraphapi.MLDeployment, error)
	CreateMLDeployment(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.MLDeployment, error)
	UpdateMLDeployment(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLDeployment, error)
	DeleteMLDeployment(ctx context.Context, orgID, id string) error
	ListMLFindings(ctx context.Context, orgID, modelID, projectID string) ([]uigraphapi.MLFinding, error)
	CreateMLFinding(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.MLFinding, error)
	UpdateMLFinding(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLFinding, error)
	DeleteMLFinding(ctx context.Context, orgID, id string) error
	ListVersionDeploymentUpdates(ctx context.Context, orgID, versionID, projectID string) ([]uigraphapi.MLVersionDeploymentUpdate, error)
	CreateVersionDeploymentUpdate(ctx context.Context, orgID, versionID string, body map[string]interface{}) (*uigraphapi.MLVersionDeploymentUpdate, error)
	SetMLModelVersionRun(ctx context.Context, orgID, versionID string, body map[string]interface{}) (*uigraphapi.MLModelVersion, error)
	LinkMLVersionEvaluations(ctx context.Context, orgID, versionID string, body map[string]interface{}) ([]uigraphapi.MLEvaluation, error)
	ListMLEvaluations(ctx context.Context, orgID string, query uigraphapi.MLEvaluationQuery) ([]uigraphapi.MLEvaluation, int, error)
	ListMLVersionEvaluations(ctx context.Context, orgID, versionID string, query uigraphapi.MLEvaluationQuery) ([]uigraphapi.MLEvaluation, int, error)
	ListMLExperimentEvaluations(ctx context.Context, orgID, experimentID string, query uigraphapi.MLEvaluationQuery) ([]uigraphapi.MLEvaluation, int, error)
	GetMLEvaluation(ctx context.Context, orgID, id string) (*uigraphapi.MLEvaluation, error)
	CreateMLEvaluation(ctx context.Context, orgID, experimentID string, body map[string]interface{}) (*uigraphapi.MLEvaluation, error)
	UpdateMLEvaluation(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.MLEvaluation, error)
	DeleteMLEvaluation(ctx context.Context, orgID, id string) error
}

type billingClient interface {
	ListCloudConnections(ctx context.Context, orgID string) ([]uigraphapi.CloudConnection, error)
	GetCloudConnection(ctx context.Context, orgID, connectionID string) (*uigraphapi.CloudConnection, error)
	CreateCloudConnection(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.CloudConnection, error)
	DeleteCloudConnection(ctx context.Context, orgID, connectionID string) error
	TestCloudConnection(ctx context.Context, orgID, connectionID string) (*uigraphapi.TestCloudConnectionResult, error)
	SyncCloudConnection(ctx context.Context, orgID, connectionID string) (bool, error)
	GetServiceCostSummary(ctx context.Context, orgID, serviceID string) (*uigraphapi.ServiceCostSummary, error)
	ListServiceCostResources(ctx context.Context, orgID, serviceID string) ([]uigraphapi.InfraResource, error)
	GetServiceCostTrend(ctx context.Context, orgID, serviceID string, days *int) ([]uigraphapi.CostTrendPoint, error)
	ListServiceCostTagRules(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceCostTagRule, error)
	CreateServiceCostTagRule(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceCostTagRule, error)
	DeleteServiceCostTagRule(ctx context.Context, orgID, serviceID, ruleID string) error
	GetResourceDailyCosts(ctx context.Context, orgID, resourceID string, days *int) ([]uigraphapi.ResourceDailyCost, error)
}

type timelineClient interface {
	ListServiceTimelineEvents(ctx context.Context, orgID, serviceID string) ([]uigraphapi.TimelineEvent, error)
	CreateTimelineEvent(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TimelineEvent, error)
	UpdateTimelineEvent(ctx context.Context, orgID, serviceID, eventID string, body map[string]interface{}) (*uigraphapi.TimelineEvent, error)
	DeleteTimelineEvent(ctx context.Context, orgID, serviceID, eventID string) error
}

type githubOnboardingClient interface {
	GetGitHubApp(ctx context.Context, orgID string) (*uigraphapi.GitHubAppInstallation, error)
	GetGitHubAppInstallURL(ctx context.Context, orgID string) (string, error)
	DisconnectGitHubApp(ctx context.Context, orgID string) error
	ListGitHubRepositories(ctx context.Context, orgID string) ([]uigraphapi.GitHubRepository, error)
	StartRepositoryOnboarding(ctx context.Context, orgID string, input uigraphapi.StartRepositoryOnboardingInput) (*uigraphapi.RepositoryOnboardingBatch, error)
	GetRepositoryOnboarding(ctx context.Context, orgID, batchID string) (*uigraphapi.RepositoryOnboardingBatch, error)
	GetLatestRepositoryOnboarding(ctx context.Context, orgID string) (*uigraphapi.RepositoryOnboardingBatch, error)
	RecheckRepositoryOnboarding(ctx context.Context, orgID, batchID, onboardingID string) (*uigraphapi.RepositoryOnboarding, error)
	RetryRepositoryOnboarding(ctx context.Context, orgID, batchID, onboardingID string) (*uigraphapi.RepositoryOnboarding, error)
}

type Resolver struct {
	Auth        authClient
	OrgAPI      orgClient
	Admin       adminClient
	AuthAPI     authProviderClient
	FolderAPI   folderClient
	DiagramAPI  diagramClient
	DocAPI      docsClient
	Component   componentClient
	UIMapAPI    uimapClient
	Catalog     catalogClient
	Dependency  dependencyClient
	Chat        chatClient
	TestPack    testPackClient
	Actor       actorClient
	CommentAPI  commentClient
	CostSavings costSavingsClient
	MLStudio    mlStudioClient
	Billing     billingClient
	AgentAPI    agentSessionClient
	Timeline    timelineClient
	GitHubAPI   githubOnboardingClient
}
