package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func OrgToModel(o *uigraphapi.Org) *model.Org {
	m := &model.Org{ID: o.ID, Name: o.Name, Disabled: o.Disabled, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt}
	if o.LogoURL != "" {
		m.LogoURL = &o.LogoURL
	}
	return m
}

func MemberToModel(m uigraphapi.Member) *model.Member {
	return &model.Member{
		UserID: m.UserID, OrgID: m.OrgID, Role: m.Role, Source: m.Source,
		Email: m.Email, Name: m.Name, TeamID: m.TeamID,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func TeamToModel(t *uigraphapi.Team) *model.Team {
	m := &model.Team{ID: t.ID, OrgID: t.OrgID, Name: t.Name, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
	if t.Email != "" {
		m.Email = &t.Email
	}
	if t.ExternalID != "" {
		m.ExternalID = &t.ExternalID
	}
	return m
}

func TeamMemberToModel(m uigraphapi.TeamMember) *model.TeamMember {
	return &model.TeamMember{TeamID: m.TeamID, UserID: m.UserID, Permission: m.Permission, CreatedAt: m.CreatedAt}
}

func ServiceAccountToModel(sa uigraphapi.ServiceAccount) *model.ServiceAccount {
	return &model.ServiceAccount{
		ID: sa.ID, OrgID: sa.OrgID, Name: sa.Name, Description: sa.Description,
		Role: sa.Role, Disabled: sa.Disabled, CreatedAt: sa.CreatedAt, UpdatedAt: sa.UpdatedAt,
	}
}

func SATokenToModel(t uigraphapi.ServiceAccountToken) *model.ServiceAccountToken {
	return &model.ServiceAccountToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name, Prefix: t.Prefix,
		ExpiresAt: t.ExpiresAt, LastUsedAt: t.LastUsedAt, Revoked: t.Revoked, CreatedAt: t.CreatedAt,
	}
}

func CreatedTokenToModel(t *uigraphapi.CreatedToken) *model.CreatedToken {
	return &model.CreatedToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name,
		Prefix: t.Prefix, Token: t.Token, CreatedAt: t.CreatedAt,
	}
}

func OrgsToModel(orgs []uigraphapi.Org) []*model.Org {
	out := make([]*model.Org, len(orgs))
	for i := range orgs {
		out[i] = OrgToModel(&orgs[i])
	}
	return out
}

func MembersToModel(members []uigraphapi.Member) []*model.Member {
	out := make([]*model.Member, len(members))
	for i, m := range members {
		out[i] = MemberToModel(m)
	}
	return out
}

func TeamsToModel(teams []uigraphapi.Team) []*model.Team {
	out := make([]*model.Team, len(teams))
	for i := range teams {
		out[i] = TeamToModel(&teams[i])
	}
	return out
}

func TeamMembersToModel(members []uigraphapi.TeamMember) []*model.TeamMember {
	out := make([]*model.TeamMember, len(members))
	for i, m := range members {
		out[i] = TeamMemberToModel(m)
	}
	return out
}

func ServiceAccountsToModel(sas []uigraphapi.ServiceAccount) []*model.ServiceAccount {
	out := make([]*model.ServiceAccount, len(sas))
	for i, sa := range sas {
		out[i] = ServiceAccountToModel(sa)
	}
	return out
}

func SATokensToModel(tokens []uigraphapi.ServiceAccountToken) []*model.ServiceAccountToken {
	out := make([]*model.ServiceAccountToken, len(tokens))
	for i, t := range tokens {
		out[i] = SATokenToModel(t)
	}
	return out
}
