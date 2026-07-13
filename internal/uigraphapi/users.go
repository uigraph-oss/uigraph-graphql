package uigraphapi

import (
	"context"
	"time"
)

type User struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Name       string     `json:"name"`
	Login      string     `json:"login"`
	Disabled   bool       `json:"disabled"`
	Role       string     `json:"role"`
	AvatarURL  string     `json:"avatarUrl,omitempty"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// ── Users (server admin) ────────────────────────────────────────────────────

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	var out struct {
		Users []User `json:"users"`
	}
	return out.Users, c.get(ctx, "/api/v1/users", &out)
}

func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	var out User
	return &out, c.get(ctx, "/api/v1/users/"+id, &out)
}

func (c *Client) CreateUser(ctx context.Context, body map[string]interface{}) (*User, error) {
	var out User
	return &out, c.post(ctx, "/api/v1/users", body, &out)
}

func (c *Client) UpdateUser(ctx context.Context, id string, body map[string]interface{}) (*User, error) {
	var out User
	return &out, c.put(ctx, "/api/v1/users/"+id, body, &out)
}

func (c *Client) DisableUser(ctx context.Context, id string) error {
	return c.del(ctx, "/api/v1/users/"+id)
}

type ServerOverview struct {
	TotalUsers  int `json:"totalUsers"`
	ActiveUsers int `json:"activeUsers"`
	TotalOrgs   int `json:"totalOrgs"`
}

func (c *Client) GetServerOverview(ctx context.Context) (*ServerOverview, error) {
	var out ServerOverview
	return &out, c.get(ctx, "/api/v1/server/overview", &out)
}

type ServerConfig struct {
	StorageBackend   string `json:"storageBackend"`
	StorageBucket    string `json:"storageBucket"`
	StorageEndpoint  string `json:"storageEndpoint"`
	VectorBackend    string `json:"vectorBackend"`
	EmbeddingBackend string `json:"embeddingBackend"`
	EmbeddingModel   string `json:"embeddingModel"`
}

func (c *Client) GetServerConfig(ctx context.Context) (*ServerConfig, error) {
	var out ServerConfig
	return &out, c.get(ctx, "/api/v1/server/config", &out)
}
