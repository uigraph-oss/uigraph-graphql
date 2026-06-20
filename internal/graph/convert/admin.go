package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func OverviewToModel(o *uigraphapi.ServerOverview) *model.ServerOverview {
	return &model.ServerOverview{
		TotalUsers: o.TotalUsers, ActiveUsers: o.ActiveUsers, TotalOrgs: o.TotalOrgs,
	}
}

func UserToModel(u *uigraphapi.User) *model.User {
	return &model.User{
		ID: u.ID, Email: u.Email, Name: u.Name, Login: u.Login,
		Disabled: u.Disabled, Role: u.Role, LastSeenAt: u.LastSeenAt,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

func UsersToModel(users []uigraphapi.User) []*model.User {
	out := make([]*model.User, len(users))
	for i := range users {
		out[i] = UserToModel(&users[i])
	}
	return out
}

func OAuthProviderToModel(p uigraphapi.OAuthProvider) *model.OAuthProvider {
	return &model.OAuthProvider{
		ID: p.ID, ProviderName: p.ProviderName, Type: p.Type, DisplayName: p.DisplayName,
		ClientID: p.ClientID, ClientSecret: p.ClientSecret,
		AuthURL: p.AuthURL, TokenURL: p.TokenURL, UserinfoURL: p.UserinfoURL, APIURL: p.APIURL,
		Scopes: p.Scopes, AllowedDomains: p.AllowedDomains, AllowSignUp: p.AllowSignUp,
		EmailClaim: p.EmailClaim, NameClaim: p.NameClaim, SubClaim: p.SubClaim,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func OAuthProvidersToModel(providers []uigraphapi.OAuthProvider) []*model.OAuthProvider {
	out := make([]*model.OAuthProvider, len(providers))
	for i := range providers {
		out[i] = OAuthProviderToModel(providers[i])
	}
	return out
}

func RoleMappingToModel(m uigraphapi.RoleMapping) *model.RoleMapping {
	return &model.RoleMapping{
		ID: m.ID, OrganizationID: m.OrganizationID,
		ClaimKey: m.ClaimKey, ClaimValue: m.ClaimValue, Role: m.Role, Scope: m.Scope,
		ResourceType: m.ResourceType, ResourceID: m.ResourceID,
	}
}

func RoleMappingsToModel(mappings []uigraphapi.RoleMapping) []*model.RoleMapping {
	out := make([]*model.RoleMapping, len(mappings))
	for i := range mappings {
		out[i] = RoleMappingToModel(mappings[i])
	}
	return out
}

func LDAPToModel(l *uigraphapi.LDAPConfig) *model.LDAPConfig {
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

func SAMLToModel(s *uigraphapi.SAMLConfig) *model.SAMLConfig {
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
