package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func MeToModel(m *uigraphapi.MeResponse) *model.Me {
	me := &model.Me{
		UserID: m.UserID, OrgID: m.OrgID,
		Email: m.Email, Name: m.Name, Login: m.Login,
		Kind: m.Kind, IsServerAdmin: m.Role == "server_admin", AuthProvider: m.AuthProvider,
	}
	if m.AvatarURL != "" {
		me.AvatarURL = &m.AvatarURL
	}
	return me
}

func OrgSummaryToModel(o uigraphapi.OrgSummary) *model.OrgSummary {
	return &model.OrgSummary{ID: o.ID, Name: o.Name, Slug: o.Slug, Role: o.Role, Active: o.Active}
}

func OrgSummariesToModel(orgs []uigraphapi.OrgSummary) []*model.OrgSummary {
	out := make([]*model.OrgSummary, len(orgs))
	for i, o := range orgs {
		out[i] = OrgSummaryToModel(o)
	}
	return out
}
