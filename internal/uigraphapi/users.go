package uigraphapi

import "context"

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
