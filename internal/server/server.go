package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/middleware"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 60 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func Run(cfg *config.Config) error {
	c := uigraphapi.New(cfg.APIBaseURL)

	resolver := &graph.Resolver{
		Auth:        c,
		OrgAPI:      c,
		Admin:       c,
		FolderAPI:   c,
		DiagramAPI:  c,
		DocAPI:      c,
		Component:   c,
		UIMapAPI:    c,
		Catalog:     c,
		Chat:        c,
		TestPack:    c,
		Actor:       c,
		CommentAPI:  c,
		CostSavings: c,
	}
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	gqlSrv := newGraphQLServer(schema, cfg.Env)

	mux := http.NewServeMux()
	mux.Handle("POST /graphql", middleware.Auth(gqlSrv))
	mux.Handle("GET /graphql", middleware.Auth(gqlSrv))

	if cfg.Env != "prod" {
		mux.Handle("GET /playground", playground.Handler("UIGraph GraphQL", "/graphql"))
		slog.Info("playground enabled", "url", "http://localhost:"+cfg.Port+"/playground")
	}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /readyz", readyzHandler(c))

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.Logging(middleware.CORS(mux)),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("uigraph-graphql listening", "port", cfg.Port, "upstream", cfg.APIBaseURL)
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	}
}

func newGraphQLServer(schema graphql.ExecutableSchema, env string) *handler.Server {
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetErrorPresenter(graph.ErrorPresenter)
	srv.Use(extension.FixedComplexityLimit(1000))
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
}

func readyzHandler(c *uigraphapi.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if err := c.Ping(ctx); err != nil {
			slog.WarnContext(ctx, "readiness check failed", "err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
