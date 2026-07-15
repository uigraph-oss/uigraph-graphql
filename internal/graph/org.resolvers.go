package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *mutationResolver) CreateOrg(ctx context.Context, input model.CreateOrgInput) (*model.Org, error) {
	o, err := r.OrgAPI.CreateOrg(ctx, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.OrgToModel(o), nil
}

func (r *mutationResolver) UpdateOrg(ctx context.Context, id string, input model.UpdateOrgInput) (*model.Org, error) {
	o, err := r.OrgAPI.UpdateOrg(ctx, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.OrgToModel(o), nil
}

func (r *mutationResolver) DeleteOrg(ctx context.Context, id string) (bool, error) {
	return true, r.OrgAPI.DeleteOrg(ctx, id)
}

func (r *mutationResolver) CompleteOnboarding(ctx context.Context, orgID string) (bool, error) {
	return true, r.OrgAPI.CompleteOnboarding(ctx, orgID)
}

func (r *mutationResolver) AddMember(ctx context.Context, orgID string, input model.AddMemberInput) (*model.Member, error) {
	m, err := r.OrgAPI.AddMember(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.MemberToModel(*m), nil
}

func (r *mutationResolver) UpdateMember(ctx context.Context, orgID string, userID string, input model.UpdateMemberInput) (*model.Member, error) {
	m, err := r.OrgAPI.UpdateMember(ctx, orgID, userID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.MemberToModel(*m), nil
}

func (r *mutationResolver) RemoveMember(ctx context.Context, orgID string, userID string) (bool, error) {
	return true, r.OrgAPI.RemoveMember(ctx, orgID, userID)
}

func (r *mutationResolver) CreateTeam(ctx context.Context, orgID string, input model.CreateTeamInput) (*model.Team, error) {
	t, err := r.OrgAPI.CreateTeam(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TeamToModel(t), nil
}

func (r *mutationResolver) UpdateTeam(ctx context.Context, orgID string, teamID string, input model.UpdateTeamInput) (*model.Team, error) {
	t, err := r.OrgAPI.UpdateTeam(ctx, orgID, teamID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TeamToModel(t), nil
}

func (r *mutationResolver) DeleteTeam(ctx context.Context, orgID string, teamID string) (bool, error) {
	return true, r.OrgAPI.DeleteTeam(ctx, orgID, teamID)
}

func (r *mutationResolver) AddTeamMember(ctx context.Context, orgID string, teamID string, userID string, permission *string) (bool, error) {
	body := map[string]interface{}{"userId": userID}
	if permission != nil {
		body["permission"] = *permission
	}
	return true, r.OrgAPI.AddTeamMember(ctx, orgID, teamID, body)
}

func (r *mutationResolver) RemoveTeamMember(ctx context.Context, orgID string, teamID string, userID string) (bool, error) {
	return true, r.OrgAPI.RemoveTeamMember(ctx, orgID, teamID, userID)
}

func (r *mutationResolver) CreateServiceAccount(ctx context.Context, orgID string, input model.CreateServiceAccountInput) (*model.ServiceAccount, error) {
	sa, err := r.OrgAPI.CreateServiceAccount(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceAccountToModel(*sa), nil
}

func (r *mutationResolver) UpdateServiceAccount(ctx context.Context, orgID string, id string, input model.UpdateServiceAccountInput) (*model.ServiceAccount, error) {
	sa, err := r.OrgAPI.UpdateServiceAccount(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceAccountToModel(*sa), nil
}

func (r *mutationResolver) DeleteServiceAccount(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.OrgAPI.DeleteServiceAccount(ctx, orgID, id)
}

func (r *mutationResolver) CreateServiceAccountToken(ctx context.Context, orgID string, saID string, input model.CreateTokenInput) (*model.CreatedToken, error) {
	t, err := r.OrgAPI.CreateServiceAccountToken(ctx, orgID, saID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.CreatedTokenToModel(t), nil
}

func (r *mutationResolver) RevokeServiceAccountToken(ctx context.Context, orgID string, saID string, tokenID string) (bool, error) {
	return true, r.OrgAPI.RevokeServiceAccountToken(ctx, orgID, saID, tokenID)
}

func (r *mutationResolver) PrepareServiceAccountAvatarUpload(ctx context.Context, orgID string, saID string) (*model.AssetUpload, error) {
	u, err := r.OrgAPI.PrepareServiceAccountAvatarUpload(ctx, orgID, saID)
	if err != nil {
		return nil, err
	}
	return &model.AssetUpload{AssetID: u.AssetID, UploadURL: u.UploadURL}, nil
}

func (r *mutationResolver) SetServiceAccountAvatar(ctx context.Context, orgID string, saID string) (bool, error) {
	return true, r.OrgAPI.SetServiceAccountAvatar(ctx, orgID, saID)
}

func (r *queryResolver) Org(ctx context.Context, id string) (*model.Org, error) {
	o, err := r.OrgAPI.GetOrg(ctx, id)
	if err != nil {
		return nil, err
	}
	return convert.OrgToModel(o), nil
}

func (r *queryResolver) Orgs(ctx context.Context) ([]*model.Org, error) {
	orgs, err := r.OrgAPI.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return convert.OrgsToModel(orgs), nil
}

func (r *queryResolver) Members(ctx context.Context, orgID string) ([]*model.Member, error) {
	members, err := r.OrgAPI.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(members))
	for i, m := range members {
		ids[i] = m.UserID
	}

	actors := map[string]*uigraphapi.Actor{}
	if len(ids) > 0 {
		actors, err = r.Resolver.Actor.ResolveActors(ctx, orgID, ids)
		if err != nil {
			return nil, err
		}
	}
	return convert.MembersToModel(members, actors), nil
}

func (r *queryResolver) Teams(ctx context.Context, orgID string) ([]*model.Team, error) {
	teams, err := r.OrgAPI.ListTeams(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return convert.TeamsToModel(teams), nil
}

func (r *queryResolver) Team(ctx context.Context, orgID string, teamID string) (*model.Team, error) {
	t, err := r.OrgAPI.GetTeam(ctx, orgID, teamID)
	if err != nil {
		return nil, err
	}
	return convert.TeamToModel(t), nil
}

func (r *queryResolver) TeamMembers(ctx context.Context, orgID string, teamID string) ([]*model.TeamMember, error) {
	members, err := r.OrgAPI.ListTeamMembers(ctx, orgID, teamID)
	if err != nil {
		return nil, err
	}
	return convert.TeamMembersToModel(members), nil
}

func (r *queryResolver) ServiceAccounts(ctx context.Context, orgID string) ([]*model.ServiceAccount, error) {
	sas, err := r.OrgAPI.ListServiceAccounts(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return convert.ServiceAccountsToModel(sas), nil
}

func (r *queryResolver) ServiceAccount(ctx context.Context, orgID string, id string) (*model.ServiceAccount, error) {
	sa, err := r.OrgAPI.GetServiceAccount(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.ServiceAccountToModel(*sa), nil
}

func (r *queryResolver) ServiceAccountTokens(ctx context.Context, orgID string, saID string) ([]*model.ServiceAccountToken, error) {
	tokens, err := r.OrgAPI.ListServiceAccountTokens(ctx, orgID, saID)
	if err != nil {
		return nil, err
	}
	return convert.SATokensToModel(tokens), nil
}

func (r *queryResolver) ServiceAccountScopes(ctx context.Context, orgID string) ([]string, error) {
	return r.OrgAPI.ListServiceAccountScopes(ctx, orgID)
}
