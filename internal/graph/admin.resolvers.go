package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *mutationResolver) CreateServerOrg(ctx context.Context, input model.CreateServerOrgInput) (*model.Org, error) {
	o, err := r.Admin.ServerCreateOrg(ctx, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.OrgToModel(o), nil
}

func (r *mutationResolver) UpdateServerOrg(ctx context.Context, id string, input model.UpdateServerOrgInput) (*model.Org, error) {
	o, err := r.Admin.ServerUpdateOrg(ctx, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.OrgToModel(o), nil
}

func (r *mutationResolver) DeleteServerOrg(ctx context.Context, id string) (bool, error) {
	return true, r.Admin.ServerDeleteOrg(ctx, id)
}

func (r *mutationResolver) PrepareServerOrgLogoUpload(ctx context.Context, orgID string) (*model.AssetUpload, error) {
	u, err := r.Admin.PrepareServerOrgLogoUpload(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return &model.AssetUpload{AssetID: u.AssetID, UploadURL: u.UploadURL}, nil
}

func (r *mutationResolver) SetServerOrgLogo(ctx context.Context, orgID string) (bool, error) {
	return true, r.Admin.SetServerOrgLogo(ctx, orgID)
}

func (r *mutationResolver) RemoveServerOrgLogo(ctx context.Context, orgID string) (bool, error) {
	return true, r.Admin.RemoveServerOrgLogo(ctx, orgID)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	u, err := r.Admin.CreateUser(ctx, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.UserToModel(u), nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
	u, err := r.Admin.UpdateUser(ctx, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.UserToModel(u), nil
}

func (r *mutationResolver) DisableUser(ctx context.Context, id string) (bool, error) {
	return true, r.Admin.DisableUser(ctx, id)
}

func (r *mutationResolver) UpsertOAuthProvider(ctx context.Context, provider string, input model.UpsertOAuthInput) (bool, error) {
	return true, r.Admin.UpsertOAuthProvider(ctx, provider, convert.ToMap(input))
}

func (r *mutationResolver) DeleteOAuthProvider(ctx context.Context, provider string) (bool, error) {
	return true, r.Admin.DeleteOAuthProvider(ctx, provider)
}

func (r *mutationResolver) PrepareOAuthProviderIconUpload(ctx context.Context, provider string) (*model.AssetUpload, error) {
	u, err := r.Admin.PrepareOAuthProviderIconUpload(ctx, provider)
	if err != nil {
		return nil, err
	}
	return &model.AssetUpload{AssetID: u.AssetID, UploadURL: u.UploadURL}, nil
}

func (r *mutationResolver) SetOAuthProviderIcon(ctx context.Context, provider string) (bool, error) {
	return true, r.Admin.SetOAuthProviderIcon(ctx, provider)
}

func (r *mutationResolver) RemoveOAuthProviderIcon(ctx context.Context, provider string) (bool, error) {
	return true, r.Admin.RemoveOAuthProviderIcon(ctx, provider)
}

func (r *mutationResolver) CreateRoleMapping(ctx context.Context, input model.CreateRoleMappingInput) (bool, error) {
	return true, r.Admin.CreateRoleMapping(ctx, convert.ToMap(input))
}

func (r *mutationResolver) DeleteRoleMapping(ctx context.Context, id string) (bool, error) {
	return true, r.Admin.DeleteRoleMapping(ctx, id)
}

func (r *mutationResolver) UpsertLdap(ctx context.Context, input model.UpsertLDAPInput) (bool, error) {
	return true, r.Admin.UpsertLDAP(ctx, convert.ToMap(input))
}

func (r *mutationResolver) DeleteLdap(ctx context.Context) (bool, error) {
	return true, r.Admin.DeleteLDAP(ctx)
}

func (r *mutationResolver) UpsertSaml(ctx context.Context, input model.UpsertSAMLInput) (bool, error) {
	return true, r.Admin.UpsertSAML(ctx, convert.ToMap(input))
}

func (r *queryResolver) ServerOverview(ctx context.Context) (*model.ServerOverview, error) {
	o, err := r.Admin.GetServerOverview(ctx)
	if err != nil {
		return nil, err
	}
	return convert.OverviewToModel(o), nil
}

func (r *queryResolver) ServerConfig(ctx context.Context) (*model.ServerConfig, error) {
	c, err := r.Admin.GetServerConfig(ctx)
	if err != nil {
		return nil, err
	}
	return convert.ServerConfigToModel(c), nil
}

func (r *queryResolver) ServerOrgs(ctx context.Context) ([]*model.Org, error) {
	orgs, err := r.Admin.ServerListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*model.Org, len(orgs))
	for i := range orgs {
		out[i] = convert.OrgToModel(&orgs[i])
	}
	return out, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.Admin.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return convert.UsersToModel(users), nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	u, err := r.Admin.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return convert.UserToModel(u), nil
}

func (r *queryResolver) OauthProviders(ctx context.Context) ([]*model.OAuthProvider, error) {
	providers, err := r.Admin.ListOAuthProviders(ctx)
	if err != nil {
		return nil, err
	}
	return convert.OAuthProvidersToModel(providers), nil
}

func (r *queryResolver) RoleMappings(ctx context.Context) ([]*model.RoleMapping, error) {
	mappings, err := r.Admin.ListRoleMappings(ctx)
	if err != nil {
		return nil, err
	}
	return convert.RoleMappingsToModel(mappings), nil
}

func (r *queryResolver) Ldap(ctx context.Context) (*model.LDAPConfig, error) {
	l, err := r.Admin.GetLDAP(ctx)
	if err != nil {
		if uigraphapi.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return convert.LDAPToModel(l), nil
}

func (r *queryResolver) Saml(ctx context.Context) (*model.SAMLConfig, error) {
	s, err := r.Admin.GetSAML(ctx)
	if err != nil {
		if uigraphapi.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return convert.SAMLToModel(s), nil
}

func (r *queryResolver) Scim(ctx context.Context) (*model.SCIMConfig, error) {
	s, err := r.Admin.GetSCIM(ctx)
	if err != nil {
		if uigraphapi.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return convert.SCIMToModel(s), nil
}
