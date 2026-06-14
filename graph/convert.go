package graph

import (
	"encoding/json"

	"github.com/uigraph/graphql/client"
	"github.com/uigraph/graphql/graph/model"
)

// toMap JSON-round-trips a struct into map[string]interface{}.
// This correctly handles optional fields: nil pointer fields are omitted
// from the resulting map (because of omitempty in the input struct tags).
func toMap(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

// rawStr returns the JSON string of a raw message, defaulting to "{}".
func rawStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "{}"
	}
	return string(b)
}

func rawArrStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "[]"
	}
	return string(b)
}

// ── Auth ─────────────────────────────────────────────────────────────────────

func meToModel(m *client.MeResponse) *model.Me {
	return &model.Me{
		UserID: m.UserID, OrgID: m.OrgID,
		Email: m.Email, Name: m.Name, Login: m.Login,
		Kind: m.Kind, Role: m.Role, AuthProvider: m.AuthProvider,
	}
}

func orgSummaryToModel(o client.OrgSummary) *model.OrgSummary {
	return &model.OrgSummary{ID: o.ID, Name: o.Name, Slug: o.Slug, Role: o.Role, Active: o.Active}
}

// ── Org ───────────────────────────────────────────────────────────────────────

func orgToModel(o *client.Org) *model.Org {
	return &model.Org{ID: o.ID, Name: o.Name, Slug: o.Slug, Disabled: o.Disabled, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt}
}

func memberToModel(m client.Member) *model.Member {
	return &model.Member{UserID: m.UserID, OrgID: m.OrgID, Role: m.Role, Source: m.Source, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func teamToModel(t *client.Team) *model.Team {
	m := &model.Team{ID: t.ID, OrgID: t.OrgID, Name: t.Name, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
	if t.Email != "" {
		m.Email = &t.Email
	}
	if t.ExternalID != "" {
		m.ExternalID = &t.ExternalID
	}
	return m
}

func teamMemberToModel(m client.TeamMember) *model.TeamMember {
	return &model.TeamMember{TeamID: m.TeamID, UserID: m.UserID, Permission: m.Permission, CreatedAt: m.CreatedAt}
}

func invitationToModel(i client.Invitation) *model.Invitation {
	return &model.Invitation{
		ID: i.ID, OrgID: i.OrgID, Email: i.Email, Role: i.Role,
		Code: i.Code, CreatedBy: i.CreatedBy, CreatedAt: i.CreatedAt, ExpiresAt: i.ExpiresAt,
	}
}

func serviceAccountToModel(sa client.ServiceAccount) *model.ServiceAccount {
	return &model.ServiceAccount{
		ID: sa.ID, OrgID: sa.OrgID, Name: sa.Name, Description: sa.Description,
		Role: sa.Role, Disabled: sa.Disabled, CreatedAt: sa.CreatedAt, UpdatedAt: sa.UpdatedAt,
	}
}

func saTokenToModel(t client.ServiceAccountToken) *model.ServiceAccountToken {
	return &model.ServiceAccountToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name, Prefix: t.Prefix,
		ExpiresAt: t.ExpiresAt, LastUsedAt: t.LastUsedAt, Revoked: t.Revoked, CreatedAt: t.CreatedAt,
	}
}

func createdTokenToModel(t *client.CreatedToken) *model.CreatedToken {
	return &model.CreatedToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name,
		Prefix: t.Prefix, Token: t.Token, CreatedAt: t.CreatedAt,
	}
}

// ── Content ───────────────────────────────────────────────────────────────────

func folderToModel(f *client.Folder) *model.Folder {
	return &model.Folder{
		ID: f.ID, OrgID: f.OrgID, ParentID: f.ParentID, Type: f.Type,
		Name: f.Name, Order: f.Order, CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func diagramToModel(d *client.Diagram) *model.Diagram {
	return &model.Diagram{
		ID: d.ID, OrgID: d.OrgID, FolderID: d.FolderID, TeamID: d.TeamID,
		Name: d.Name, ContentKey: d.ContentKey, ContentHash: d.ContentHash,
		PreviewImageFileID: d.PreviewImageFileID, PreviewImageURL: d.PreviewImageURL,
		Source: d.Source, CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func diagramVersionToModel(v client.DiagramVersion) *model.DiagramVersion {
	return &model.DiagramVersion{
		ID: v.ID, DiagramID: v.DiagramID, VersionNumber: v.VersionNumber,
		Label: v.Label, ContentKey: v.ContentKey, ContentHash: v.ContentHash,
		IsAutoVersion: v.IsAutoVersion, Source: v.Source, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func flowComponentFieldToModel(f client.FlowDiagramComponentField) *model.FlowDiagramComponentField {
	return &model.FlowDiagramComponentField{
		FlowDiagramComponentFieldID: f.FlowDiagramComponentFieldID,
		Label:                       f.Label,
		Type:                        f.Type,
		Required:                    f.Required,
		Readonly:                    f.Readonly,
		Options:                     f.Options,
		Order:                       f.Order,
	}
}

func flowComponentToModel(c client.FlowDiagramComponent) *model.FlowDiagramComponent {
	fields := make([]*model.FlowDiagramComponentField, len(c.FlowDiagramComponentFields))
	for i, f := range c.FlowDiagramComponentFields {
		fields[i] = flowComponentFieldToModel(f)
	}
	return &model.FlowDiagramComponent{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, OrganizationID: c.OrganizationID,
		FlowDiagramComponentFields: fields,
	}
}

func flowComponentsToModel(components []client.FlowDiagramComponent) []*model.FlowDiagramComponent {
	out := make([]*model.FlowDiagramComponent, len(components))
	for i, c := range components {
		out[i] = flowComponentToModel(c)
	}
	return out
}

func diagramImageToModel(img client.DiagramImage) *model.DiagramImage {
	return &model.DiagramImage{
		DiagramImageID: img.DiagramImageID, DiagramID: img.DiagramID,
		FileID: img.FileID, FileURL: img.FileURL, FileName: img.FileName,
		Order: img.Order, CreatedBy: img.CreatedBy, CreatedAt: img.CreatedAt,
	}
}

func diagramImagesToModel(images []client.DiagramImage) []*model.DiagramImage {
	out := make([]*model.DiagramImage, len(images))
	for i, img := range images {
		out[i] = diagramImageToModel(img)
	}
	return out
}

func uimapToModel(m *client.UIMap) *model.UIMap {
	return &model.UIMap{
		ID: m.ID, OrgID: m.OrgID, FolderID: m.FolderID, TeamID: m.TeamID,
		Name: m.Name, Description: m.Description, Status: m.Status,
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func frameToModel(f *client.Frame) *model.Frame {
	return &model.Frame{
		ID: f.ID, MapID: f.MapID, OrgID: f.OrgID, ParentFrameID: f.ParentFrameID,
		Name: f.Name, Description: f.Description, TemplateType: f.TemplateType,
		ScreenshotKey: f.ScreenshotKey, ScreenshotContentHash: f.ScreenshotContentHash,
		ScreenshotURL: f.ScreenshotURL,
		Status: f.Status, Order: f.Order, Source: f.Source,
		CreatedBy: f.CreatedBy, UpdatedBy: f.UpdatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func focalPointToModel(fp *client.FocalPoint) *model.FocalPoint {
	return &model.FocalPoint{
		ID: fp.ID, FrameID: fp.FrameID, OrgID: fp.OrgID,
		Name: fp.Name, LocationX: fp.LocationX, LocationY: fp.LocationY,
		Visibility: fp.Visibility, IsActive: fp.IsActive,
		CreatedBy: fp.CreatedBy, UpdatedBy: fp.UpdatedBy, CreatedAt: fp.CreatedAt, UpdatedAt: fp.UpdatedAt,
	}
}

func canvasToModel(c *client.Canvas) *model.Canvas {
	return &model.Canvas{
		MapID: c.MapID, OrgID: c.OrgID,
		Zoom: c.Zoom, NavigationX: c.NavigationX, NavigationY: c.NavigationY,
		FramePositions: rawStr(c.FramePositions),
		UpdatedAt: c.UpdatedAt,
	}
}

func frameGroupToModel(g *client.FrameGroup) *model.FrameGroup {
	return &model.FrameGroup{
		ID: g.ID, FrameID: g.FrameID, OrgID: g.OrgID,
		Name: g.Name, Description: g.Description,
		LocationX: g.LocationX, LocationY: g.LocationY,
		Width: g.Width, Height: g.Height, Order: g.Order, IsActive: g.IsActive,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy,
		CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func frameGroupsToModel(gs []client.FrameGroup) []*model.FrameGroup {
	out := make([]*model.FrameGroup, len(gs))
	for i := range gs {
		out[i] = frameGroupToModel(&gs[i])
	}
	return out
}

func frameLinkToModel(l *client.FrameLink) *model.FrameLink {
	return &model.FrameLink{
		ID: l.ID, FrameID: l.FrameID, OrgID: l.OrgID, Kind: l.Kind,
		TargetFrameID: l.TargetFrameID, TargetMapID: l.TargetMapID,
		Label: l.Label, LocationX: l.LocationX, LocationY: l.LocationY, IsActive: l.IsActive,
		CreatedBy: l.CreatedBy, UpdatedBy: l.UpdatedBy,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func frameLinksToModel(ls []client.FrameLink) []*model.FrameLink {
	out := make([]*model.FrameLink, len(ls))
	for i := range ls {
		out[i] = frameLinkToModel(&ls[i])
	}
	return out
}

func focalPointMetaToModel(m *client.FocalPointMeta) *model.FocalPointMeta {
	return &model.FocalPointMeta{
		ID: m.ID, FocalPointID: m.FocalPointID, OrgID: m.OrgID, FrameID: m.FrameID,
		ComponentID: m.ComponentID, ComponentLinkID: m.ComponentLinkID,
		ComponentImages:      rawArrStr(m.ComponentImages),
		ComponentFlowDiagram: m.ComponentFlowDiagram,
		ComponentModalFields: rawArrStr(m.ComponentModalFields),
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func focalPointMetasToModel(ms []client.FocalPointMeta) []*model.FocalPointMeta {
	out := make([]*model.FocalPointMeta, len(ms))
	for i := range ms {
		out[i] = focalPointMetaToModel(&ms[i])
	}
	return out
}

// focalPointMetaBody decodes the JSON-string fields (componentImages,
// componentModalFields) into parsed JSON so the API stores real JSON arrays
// rather than quoted strings.
func focalPointMetaBody(body map[string]interface{}) map[string]interface{} {
	for _, key := range []string{"componentImages", "componentModalFields"} {
		if s, ok := body[key].(string); ok {
			var raw interface{}
			if err := unmarshalJSONString(s, &raw); err == nil {
				body[key] = raw
			}
		}
	}
	return body
}

// ── Catalog ───────────────────────────────────────────────────────────────────

func serviceToModel(s *client.Service) *model.Service {
	return &model.Service{
		ID: s.ID, OrgID: s.OrgID, FolderID: s.FolderID, TeamID: s.TeamID,
		Name: s.Name, Slug: s.Slug, Description: s.Description,
		Status: s.Status, Tier: s.Tier, Category: s.Category, Language: s.Language,
		GitRepoURL: s.GitRepoURL, JiraProjectURL: s.JiraProjectURL,
		SlackChannelURL: s.SlackChannelURL, LastCommitSha: s.LastCommitSha,
		Labels:    s.Labels,
		Metadata:  rawStr(s.Metadata),
		CreatedBy: s.CreatedBy, UpdatedBy: s.UpdatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func apiGroupToModel(g *client.APIGroup) *model.APIGroup {
	return &model.APIGroup{
		ID: g.ID, ServiceID: g.ServiceID, OrgID: g.OrgID,
		Name: g.Name, Version: g.Version, Label: g.Label, Protocol: g.Protocol,
		SpecKey: g.SpecKey, SpecHash: g.SpecHash,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func apiGroupVersionToModel(v client.APIGroupVersion) *model.APIGroupVersion {
	return &model.APIGroupVersion{
		ID: v.ID, APIGroupID: v.APIGroupID, VersionNumber: v.VersionNumber,
		Label: v.Label, SpecKey: v.SpecKey, SpecHash: v.SpecHash,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func apiEndpointToModel(e *client.APIEndpoint) *model.APIEndpoint {
	return &model.APIEndpoint{
		ID: e.ID, APIGroupID: e.APIGroupID, ServiceID: e.ServiceID, OrgID: e.OrgID,
		OperationID: e.OperationID, Method: e.Method, Path: e.Path,
		Summary: e.Summary, Description: e.Description, Tags: e.Tags,
		Parameters:  rawArrStr(e.Parameters),
		RequestBody: rawStr(e.RequestBody),
		Responses:   rawStr(e.Responses),
		Order:       e.Order,
		CreatedBy: e.CreatedBy, UpdatedBy: e.UpdatedBy, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

// ── List helpers ──────────────────────────────────────────────────────────────

func orgsToModel(orgs []client.Org) []*model.Org {
	out := make([]*model.Org, len(orgs))
	for i := range orgs {
		out[i] = orgToModel(&orgs[i])
	}
	return out
}

func orgSummariesToModel(orgs []client.OrgSummary) []*model.OrgSummary {
	out := make([]*model.OrgSummary, len(orgs))
	for i, o := range orgs {
		out[i] = orgSummaryToModel(o)
	}
	return out
}

func membersToModel(members []client.Member) []*model.Member {
	out := make([]*model.Member, len(members))
	for i, m := range members {
		out[i] = memberToModel(m)
	}
	return out
}

func teamsToModel(teams []client.Team) []*model.Team {
	out := make([]*model.Team, len(teams))
	for i := range teams {
		out[i] = teamToModel(&teams[i])
	}
	return out
}

func teamMembersToModel(members []client.TeamMember) []*model.TeamMember {
	out := make([]*model.TeamMember, len(members))
	for i, m := range members {
		out[i] = teamMemberToModel(m)
	}
	return out
}

func invitationsToModel(invs []client.Invitation) []*model.Invitation {
	out := make([]*model.Invitation, len(invs))
	for i, inv := range invs {
		out[i] = invitationToModel(inv)
	}
	return out
}

func serviceAccountsToModel(sas []client.ServiceAccount) []*model.ServiceAccount {
	out := make([]*model.ServiceAccount, len(sas))
	for i, sa := range sas {
		out[i] = serviceAccountToModel(sa)
	}
	return out
}

func saTokensToModel(tokens []client.ServiceAccountToken) []*model.ServiceAccountToken {
	out := make([]*model.ServiceAccountToken, len(tokens))
	for i, t := range tokens {
		out[i] = saTokenToModel(t)
	}
	return out
}

func foldersToModel(folders []client.Folder) []*model.Folder {
	out := make([]*model.Folder, len(folders))
	for i := range folders {
		out[i] = folderToModel(&folders[i])
	}
	return out
}

func diagramsToModel(diagrams []client.Diagram) []*model.Diagram {
	out := make([]*model.Diagram, len(diagrams))
	for i := range diagrams {
		out[i] = diagramToModel(&diagrams[i])
	}
	return out
}

func diagramVersionsToModel(versions []client.DiagramVersion) []*model.DiagramVersion {
	out := make([]*model.DiagramVersion, len(versions))
	for i, v := range versions {
		out[i] = diagramVersionToModel(v)
	}
	return out
}

func uimapsToModel(maps []client.UIMap) []*model.UIMap {
	out := make([]*model.UIMap, len(maps))
	for i := range maps {
		out[i] = uimapToModel(&maps[i])
	}
	return out
}

func framesToModel(frames []client.Frame) []*model.Frame {
	out := make([]*model.Frame, len(frames))
	for i := range frames {
		out[i] = frameToModel(&frames[i])
	}
	return out
}

func focalPointsToModel(fps []client.FocalPoint) []*model.FocalPoint {
	out := make([]*model.FocalPoint, len(fps))
	for i := range fps {
		out[i] = focalPointToModel(&fps[i])
	}
	return out
}

func servicesToModel(services []client.Service) []*model.Service {
	out := make([]*model.Service, len(services))
	for i := range services {
		out[i] = serviceToModel(&services[i])
	}
	return out
}

func apiGroupsToModel(groups []client.APIGroup) []*model.APIGroup {
	out := make([]*model.APIGroup, len(groups))
	for i := range groups {
		out[i] = apiGroupToModel(&groups[i])
	}
	return out
}

func apiGroupVersionsToModel(versions []client.APIGroupVersion) []*model.APIGroupVersion {
	out := make([]*model.APIGroupVersion, len(versions))
	for i, v := range versions {
		out[i] = apiGroupVersionToModel(v)
	}
	return out
}

func apiEndpointsToModel(endpoints []client.APIEndpoint) []*model.APIEndpoint {
	out := make([]*model.APIEndpoint, len(endpoints))
	for i := range endpoints {
		out[i] = apiEndpointToModel(&endpoints[i])
	}
	return out
}

// ── Users ─────────────────────────────────────────────────────────────────────

func userToModel(u *client.User) *model.User {
	return &model.User{
		ID: u.ID, Email: u.Email, Name: u.Name, Login: u.Login,
		Disabled: u.Disabled, Role: u.Role, LastSeenAt: u.LastSeenAt,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

func usersToModel(users []client.User) []*model.User {
	out := make([]*model.User, len(users))
	for i := range users {
		out[i] = userToModel(&users[i])
	}
	return out
}

// ── SSO ───────────────────────────────────────────────────────────────────────

func oauthProviderToModel(p client.OAuthProvider) *model.OAuthProvider {
	return &model.OAuthProvider{
		ID: p.ID, ProviderName: p.ProviderName, Type: p.Type, DisplayName: p.DisplayName,
		ClientID: p.ClientID, ClientSecret: p.ClientSecret,
		AuthURL: p.AuthURL, TokenURL: p.TokenURL, UserinfoURL: p.UserinfoURL, APIURL: p.APIURL,
		Scopes: p.Scopes, AllowedDomains: p.AllowedDomains, AllowSignUp: p.AllowSignUp,
		EmailClaim: p.EmailClaim, NameClaim: p.NameClaim, SubClaim: p.SubClaim,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func oauthProvidersToModel(providers []client.OAuthProvider) []*model.OAuthProvider {
	out := make([]*model.OAuthProvider, len(providers))
	for i := range providers {
		out[i] = oauthProviderToModel(providers[i])
	}
	return out
}

func roleMappingToModel(m client.RoleMapping) *model.RoleMapping {
	return &model.RoleMapping{
		ID: m.ID, OrganizationID: m.OrganizationID,
		ClaimKey: m.ClaimKey, ClaimValue: m.ClaimValue, Role: m.Role, Scope: m.Scope,
		ResourceType: m.ResourceType, ResourceID: m.ResourceID,
	}
}

func roleMappingsToModel(mappings []client.RoleMapping) []*model.RoleMapping {
	out := make([]*model.RoleMapping, len(mappings))
	for i := range mappings {
		out[i] = roleMappingToModel(mappings[i])
	}
	return out
}

func ldapToModel(l *client.LDAPConfig) *model.LDAPConfig {
	return &model.LDAPConfig{
		ID: l.ID, Host: l.Host, Port: l.Port,
		UseSsl: l.UseSSL, StartTLS: l.StartTLS, SkipTLSVerify: l.SkipTLSVerify,
		BindDn: l.BindDN, BindPassword: l.BindPassword,
		SearchBaseDn: l.SearchBaseDN, SearchFilter: l.SearchFilter,
		EmailAttribute: l.EmailAttribute, NameAttribute: l.NameAttribute,
		UsernameAttribute: l.UsernameAttribute, MemberOfAttribute: l.MemberOfAttribute,
		AllowSignUp: l.AllowSignUp, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func samlToModel(s *client.SAMLConfig) *model.SAMLConfig {
	return &model.SAMLConfig{
		ID: s.ID, IdpMetadataURL: s.IDPMetadataURL, IdpMetadataXML: s.IDPMetadataXML,
		IdpEntityID: s.IDPEntityID, IdpSsoURL: s.IDPSsoURL, IdpCert: s.IDPCert,
		SpEntityID: s.SPEntityID, SpCert: s.SPCert, SpKey: s.SPKey,
		SignRequests: s.SignRequests, NameIDFormat: s.NameIDFormat,
		EmailAttribute: s.EmailAttribute, NameAttribute: s.NameAttribute, LoginAttribute: s.LoginAttribute,
		GroupsAttribute: s.GroupsAttribute, AllowSignUp: s.AllowSignUp,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}
