package main

import (
	"fmt"
	"log"
	"myHabr/internal/config"
	"myHabr/internal/graph"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultPort = "8080"

func main() {
	cfg := config.MustLoad()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	grcpConnUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.GRPC.UserContainerIP, cfg.GRPC.UserPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConnUser.Close()

	grcpConnPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.GRPC.PostContainerIP, cfg.GRPC.PostPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConnPost.Close()

	resolver := graph.NewResolver(grcpConnUser, grcpConnPost)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
