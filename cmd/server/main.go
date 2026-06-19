package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/uigraph/graphql/internal/uigraphapi"
	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/middleware"
)

func main() {
	cfg := config.Load()

	c := uigraphapi.New(cfg.APIBaseURL)

	resolver := &graph.Resolver{
		Auth:      c,
		OrgAPI:    c,
		Admin:     c,
		FolderAPI: c,
		DiagramAPI: c,
		Component: c,
		UIMap:     c,
		Catalog:   c,
		TestPack:  c,
		Actor:     c,
	}
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	srv := newServer(schema, cfg.Env)

	mux := http.NewServeMux()
	mux.Handle("POST /graphql", middleware.Auth(srv))
	mux.Handle("GET /graphql", middleware.Auth(srv)) // GET for introspection tools

	if cfg.Env != "prod" {
		mux.Handle("GET /playground", playground.Handler("UIGraph GraphQL", "/graphql"))
		log.Printf("GraphQL Playground at http://localhost:%s/playground", cfg.Port)
	}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("uigraph-graphql listening on :%s → %s", cfg.Port, cfg.APIBaseURL)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}

func newServer(schema graphql.ExecutableSchema, env string) *handler.Server {
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
}
