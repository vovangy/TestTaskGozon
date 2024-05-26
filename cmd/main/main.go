package main

import (
	"context"
	"fmt"
	"log"
	"myHabr/internal/config"
	"myHabr/internal/graph"
	"myHabr/internal/middleware/jwt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"
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

	c := graph.Config{Resolvers: resolver}
	c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		httpRequest := graphql.GetRequestContext(ctx)
		log.Println(httpRequest)
		cookie := httpRequest.Headers.Get("Authorization")
		if cookie == "" {
			return nil, gqlerror.Errorf("Authorization header required")
		}

		claims, err := jwt.ParseToken(cookie)
		if err != nil {
			return nil, gqlerror.Errorf("Invalid token")
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			return nil, gqlerror.Errorf("Invalid token")
		}

		if timeExp.Before(time.Now()) {
			return nil, gqlerror.Errorf("Invalid token")
		}

		id, err := jwt.ParseClaims(claims)
		if err != nil {
			return nil, gqlerror.Errorf("Invalid token")
		}

		return next(context.WithValue(ctx, "userid", id))
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "__http_request", r)
		srv.ServeHTTP(w, r.WithContext(ctx))
	}))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
