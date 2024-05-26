package graph

import (
	"context"
	"net/http"
	"time"

	"myHabr/internal/middleware/jwt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const CookieName = "jwt-myHabr"

// authDirective is the directive for authenticating users
func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver, _ interface{}) (interface{}, error) {
	httpRequest := graphql.GetOperationContext(ctx).Variables["__http_request"].(*http.Request)
	cookie := httpRequest.Header.Get("Authorization")
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
