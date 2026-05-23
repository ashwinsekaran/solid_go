package main

import (
	"log"
	"net/http"
	"os"
	graph2 "solid_go/db/graphql/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// handler.New wires the generated schema (types + resolvers) into a gqlgen server.
	// graph.Config lets you inject middleware (field-level auth, logging, etc.) if needed.
	srv := handler.New(graph2.NewExecutableSchema(graph2.Config{Resolvers: &graph2.Resolver{}}))

	// Transports determine which HTTP methods the server accepts for GraphQL operations.
	srv.AddTransport(transport.Options{}) // responds to OPTIONS for CORS preflight
	srv.AddTransport(transport.GET{})     // allows queries via GET + query string
	srv.AddTransport(transport.POST{})    // standard: queries/mutations via POST body

	// LRU query cache: parsed query documents are expensive to re-parse on every request.
	// Caching up to 1000 distinct queries avoids redundant AST work.
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Introspection lets tools like GraphQL Playground discover the schema at runtime.
	srv.Use(extension.Introspection{})

	// Automatic Persisted Queries (APQ): clients send a hash of the query instead of
	// the full query string. On a cache hit, the server runs the query without parsing;
	// on a miss it asks the client to resend the full query, which then gets cached.
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// "/" serves the interactive GraphQL Playground — a browser-based IDE for exploring the API.
	// "/query" is the actual GraphQL endpoint that clients POST to.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
