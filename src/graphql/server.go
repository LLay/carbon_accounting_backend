package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/michelaquino/golang_api_skeleton/src/graphql/graph"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := graph.NewResolver()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: r,
	}))

	// Create a CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3001", "http://google.com:1234"}, // Add your frontend's URL here
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},                          // Add the allowed HTTP methods
		AllowedHeaders: []string{"*"},                                               // Add the allowed headers here
	})
	// Wrap your GraphQL handler with CORS middleware
	http.Handle("/query", c.Handler(srv)) // Wrap your GraphQL handler with CORS middleware

	// Handle CORS preflight requests (HTTP OPTIONS) for the '/query' endpoint
	// handlerWithCORS := c.Handler(srv)
	// http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("/query Request: %+v\n", r)
	// 	if r.Method == "OPTIONS" {
	// 		fmt.Printf("OPTIONS /query")
	// 		// Respond to CORS preflight requests
	// 		w.WriteHeader(http.StatusOK)
	// 		return
	// 	}
	// 	// For non-preflight requests, pass the request to the GraphQL handler
	// 	handlerWithCORS.ServeHTTP(w, r)
	// })

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// http.Handle("/query", handler.GraphQL(
	// 	resolvers.NewExecutableSchema(resolvers.Config{Resolvers: &resolvers.Resolver{}}),
	// 	handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
	// 		ctx = context.WithValue(ctx, "key", "value") // Add context values if needed
	// 		return next(ctx)
	// 	}),
	// 	handler.ErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
	// 		return graphql.DefaultErrorPresenter(ctx, e)
	// 	}),
	// ))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
