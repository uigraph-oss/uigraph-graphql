package client

import (
	"context"
	"fmt"
)

// ── Orgs ──────────────────────────────────────────────────────────────────────

func (c *Client) ListOrgs(ctx context.Context) ([]Org, error) {
	var out struct {
		Orgs []Org `json:"orgs"`
	}
	return out.Orgs, c.get(ctx, "/api/v1/orgs", &out)
}

func (c *Client) GetOrg(ctx context.Context, id string) (*Org, error) {
	var out Org
	return &out, c.get(ctx, "/api/v1/orgs/"+id, &out)
}

func (c *Client) CreateOrg(ctx context.Context, body map[string]interface{}) (*Org, error) {
	var out Org
	return &out, c.post(ctx, "/api/v1/orgs", body, &out)
}

func (c *Client) UpdateOrg(ctx context.Context, id string, body map[string]interface{}) (*Org, error) {
	var out Org
	return &out, c.put(ctx, "/api/v1/orgs/"+id, body, &out)
}

func (c *Client) DeleteOrg(ctx context.Context, id string) error {
	return c.del(ctx, "/api/v1/orgs/"+id)
}

// ── Members ───────────────────────────────────────────────────────────────────

func (c *Client) ListMembers(ctx context.Context, orgID string) ([]Member, error) {
	var out struct {
		Members []Member `json:"members"`
	}
	return out.Members, c.get(ctx, "/api/v1/orgs/"+orgID+"/members", &out)
}

func (c *Client) AddMember(ctx context.Context, orgID string, body map[string]interface{}) (*Member, error) {
	var out Member
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/members", body, &out)
}

func (c *Client) UpdateMemberRole(ctx context.Context, orgID, userID string, body map[string]interface{}) (*Member, error) {
	var out Member
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/members/%s", orgID, userID), body, &out)
}

func (c *Client) RemoveMember(ctx context.Context, orgID, userID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/members/%s", orgID, userID))
}

// ── Teams ─────────────────────────────────────────────────────────────────────

func (c *Client) ListTeams(ctx context.Context, orgID string) ([]Team, error) {
	var out struct {
		Teams []Team `json:"teams"`
	}
	return out.Teams, c.get(ctx, "/api/v1/orgs/"+orgID+"/teams", &out)
}

func (c *Client) GetTeam(ctx context.Context, orgID, teamID string) (*Team, error) {
	var out Team
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/teams/%s", orgID, teamID), &out)
}

func (c *Client) CreateTeam(ctx context.Context, orgID string, body map[string]interface{}) (*Team, error) {
	var out Team
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/teams", body, &out)
}

func (c *Client) UpdateTeam(ctx context.Context, orgID, teamID string, body map[string]interface{}) (*Team, error) {
	var out Team
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/teams/%s", orgID, teamID), body, &out)
}

func (c *Client) DeleteTeam(ctx context.Context, orgID, teamID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/teams/%s", orgID, teamID))
}

func (c *Client) ListTeamMembers(ctx context.Context, orgID, teamID string) ([]TeamMember, error) {
	var out struct {
		Members []TeamMember `json:"members"`
	}
	return out.Members, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/teams/%s/members", orgID, teamID), &out)
}

func (c *Client) AddTeamMember(ctx context.Context, orgID, teamID string, body map[string]interface{}) error {
	return c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/teams/%s/members", orgID, teamID), body, nil)
}

func (c *Client) RemoveTeamMember(ctx context.Context, orgID, teamID, userID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/teams/%s/members/%s", orgID, teamID, userID))
}

// ── Invitations ───────────────────────────────────────────────────────────────

func (c *Client) ListInvitations(ctx context.Context, orgID string) ([]Invitation, error) {
	var out struct {
		Invitations []Invitation `json:"invitations"`
	}
	return out.Invitations, c.get(ctx, "/api/v1/orgs/"+orgID+"/invitations", &out)
}

func (c *Client) CreateInvitation(ctx context.Context, orgID string, body map[string]interface{}) (*Invitation, error) {
	var out Invitation
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/invitations", body, &out)
}

func (c *Client) RevokeInvitation(ctx context.Context, orgID, invitationID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/invitations/%s", orgID, invitationID))
}

// ── Service Accounts ──────────────────────────────────────────────────────────

func (c *Client) ListServiceAccounts(ctx context.Context, orgID string) ([]ServiceAccount, error) {
	var out struct {
		ServiceAccounts []ServiceAccount `json:"serviceAccounts"`
	}
	return out.ServiceAccounts, c.get(ctx, "/api/v1/orgs/"+orgID+"/service-accounts", &out)
}

func (c *Client) GetServiceAccount(ctx context.Context, orgID, id string) (*ServiceAccount, error) {
	var out ServiceAccount
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/service-accounts/%s", orgID, id), &out)
}

func (c *Client) CreateServiceAccount(ctx context.Context, orgID string, body map[string]interface{}) (*ServiceAccount, error) {
	var out ServiceAccount
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/service-accounts", body, &out)
}

func (c *Client) UpdateServiceAccount(ctx context.Context, orgID, id string, body map[string]interface{}) (*ServiceAccount, error) {
	var out ServiceAccount
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/service-accounts/%s", orgID, id), body, &out)
}

func (c *Client) DeleteServiceAccount(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/service-accounts/%s", orgID, id))
}

func (c *Client) ListServiceAccountTokens(ctx context.Context, orgID, saID string) ([]ServiceAccountToken, error) {
	var out struct {
		Tokens []ServiceAccountToken `json:"tokens"`
	}
	return out.Tokens, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/service-accounts/%s/tokens", orgID, saID), &out)
}

func (c *Client) CreateServiceAccountToken(ctx context.Context, orgID, saID string, body map[string]interface{}) (*CreatedToken, error) {
	var out CreatedToken
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/service-accounts/%s/tokens", orgID, saID), body, &out)
}

func (c *Client) RevokeServiceAccountToken(ctx context.Context, orgID, saID, tokenID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/service-accounts/%s/tokens/%s", orgID, saID, tokenID))
}
