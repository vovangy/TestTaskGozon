package graph

import (
	"myHabr/internal/models"
	genPost "myHabr/internal/posts/delivery/grpc/gen"
	genUser "myHabr/internal/users/delivery/grpc/gen"

	"google.golang.org/grpc"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var users = []*models.User{{ID: "1", Username: "Vova"}}
var posts = []*models.Post{{ID: "1", Title: "Aba", Content: "CHTOTO", Author: users[0], Comments: nil}}

// SignUp is the resolver for the signUp field.
/*func (r *mutationResolver) SignUp(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: SignUp - signUp"))
}*/

type Resolver struct {
	grpcUserClient genUser.UserClient
	grpcPostClient genPost.PostClient
}

func NewResolver(grpcUserClient *grpc.ClientConn, grpcPostClient *grpc.ClientConn) *Resolver {
	return &Resolver{grpcUserClient: genUser.NewUserClient(grpcUserClient), grpcPostClient: genPost.NewPostClient(grpcPostClient)}
}
