package graph

//go:generate go run github.com/99designs/gqlgen generate

import "github.com/uigraph/graphql/client"

// Resolver is the root dependency-injection struct for all resolvers.
type Resolver struct {
	Client *client.Client
}
