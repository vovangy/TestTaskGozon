package graph

import (
	"context"
	"fmt"
	"myHabr/internal/graph/model"
	"myHabr/internal/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var users = []*models.User{{ID: "1", Username: "Vova", Posts: nil, Comments: nil}}
var posts = []*models.Post{{ID: "1", Title: "Aba", Content: "CHTOTO", Author: users[0], Comments: nil}}

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: SignUp - signUp"))
}

type Resolver struct{}
