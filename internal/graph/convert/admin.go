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

func ServerConfigToModel(c *uigraphapi.ServerConfig) *model.ServerConfig {
	return &model.ServerConfig{
		StorageBackend: c.StorageBackend, StorageBucket: c.StorageBucket,
		StorageEndpoint: c.StorageEndpoint, VectorBackend: c.VectorBackend,
		EmbeddingBackend: c.EmbeddingBackend, EmbeddingModel: c.EmbeddingModel,
	}
}

func UserToModel(u *uigraphapi.User) *model.User {
	m := &model.User{
		ID: u.ID, Email: u.Email, Name: u.Name, Login: u.Login,
		Disabled: u.Disabled, Role: u.Role, LastSeenAt: u.LastSeenAt,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
	if u.AvatarURL != "" {
		m.AvatarURL = &u.AvatarURL
	}
	return m
}

func UsersToModel(users []uigraphapi.User) []*model.User {
	out := make([]*model.User, len(users))
	for i := range users {
		out[i] = UserToModel(&users[i])
	}
	return out
}

func SCIMToModel(s *uigraphapi.SCIMConfig) *model.SCIMConfig {
	return &model.SCIMConfig{ID: s.ID}
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
