package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func AuthProviderToModel(p *uigraphapi.AuthProvider) *model.AuthProvider {
	return &model.AuthProvider{
		ID: p.ID, Slug: p.Slug, OrgID: p.OrgID, Kind: p.Kind, Type: p.Type,
		DisplayName: p.DisplayName, IconURL: p.IconURL,
		Enabled: p.Enabled, AllowSignUp: p.AllowSignUp,
		AllowedDomains: p.AllowedDomains, DefaultRole: p.DefaultRole,

		ClientID: p.ClientID, ClientSecret: p.ClientSecret,
		AuthURL: p.AuthURL, TokenURL: p.TokenURL, UserinfoURL: p.UserinfoURL, APIURL: p.APIURL,
		Scopes:     p.Scopes,
		EmailClaim: p.EmailClaim, NameClaim: p.NameClaim, SubClaim: p.SubClaim, GroupsClaim: p.GroupsClaim,

		IdpMetadataURL: p.IDPMetadataURL, IdpMetadataXML: p.IDPMetadataXML,
		IdpEntityID: p.IDPEntityID, IdpSsoURL: p.IDPSsoURL, IdpCert: p.IDPCert,
		SpEntityID: p.SPEntityID, SpCert: p.SPCert, SpKey: p.SPKey,
		SignRequests: p.SignRequests, NameIDFormat: p.NameIDFormat,
		EmailAttribute: p.EmailAttribute, NameAttribute: p.NameAttribute,
		GroupsAttribute: p.GroupsAttribute,

		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func AuthProvidersToModel(providers []uigraphapi.AuthProvider) []*model.AuthProvider {
	out := make([]*model.AuthProvider, len(providers))
	for i := range providers {
		out[i] = AuthProviderToModel(&providers[i])
	}
	return out
}

func SamlSpMetadataToModel(m *uigraphapi.SamlSpMetadata) *model.SamlSpMetadata {
	return &model.SamlSpMetadata{
		EntityID: m.EntityID, AcsURL: m.ACSURL,
		MetadataURL: m.MetadataURL, Metadata: m.Metadata,
	}
}

func AuthRoleMappingToModel(m *uigraphapi.AuthRoleMapping) *model.AuthRoleMapping {
	return &model.AuthRoleMapping{
		ID: m.ID, ProviderID: m.ProviderID, Priority: m.Priority,
		AttributeKey: m.AttributeKey, Operator: m.Operator,
		AttributeValue: m.AttributeValue, Role: m.Role,
	}
}

func AuthRoleMappingsToModel(mappings []uigraphapi.AuthRoleMapping) []*model.AuthRoleMapping {
	out := make([]*model.AuthRoleMapping, len(mappings))
	for i := range mappings {
		out[i] = AuthRoleMappingToModel(&mappings[i])
	}
	return out
}

func OrgDomainToModel(d *uigraphapi.OrgDomain) *model.OrgDomain {
	return &model.OrgDomain{
		ID: d.ID, OrgID: d.OrgID, Domain: d.Domain,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func OrgDomainsToModel(domains []uigraphapi.OrgDomain) []*model.OrgDomain {
	out := make([]*model.OrgDomain, len(domains))
	for i := range domains {
		out[i] = OrgDomainToModel(&domains[i])
	}
	return out
}

func MappingOperatorsToModel(ops []uigraphapi.MappingOperator) []*model.MappingOperator {
	out := make([]*model.MappingOperator, len(ops))
	for i := range ops {
		out[i] = &model.MappingOperator{Name: ops[i].Name, TakesValue: ops[i].TakesValue}
	}
	return out
}
