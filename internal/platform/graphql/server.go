package graphql

import (
	"net/http"

	gqlcore "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// NewHandler returns an http.Handler that serves the GraphQL endpoint at /query
// and, when enabled, the GraphQL Playground at /.
func NewHandler(schema gqlcore.ExecutableSchema, playgroundEnabled bool) http.Handler {
	srv := handler.NewDefaultServer(schema)

	mux := http.NewServeMux()
	mux.Handle("/query", srv)
	if playgroundEnabled {
		mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	}

	return mux
}
