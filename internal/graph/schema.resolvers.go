package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"
	"log/slog"
	"myHabr/internal/graph/model"
	"myHabr/internal/models"
	"myHabr/internal/users/delivery/grpc/gen"
)

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	response, err := r.grpcUserClient.SignUp(ctx, &gen.SignInUpRequest{Username: input.Username, Password: input.Password})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("User created GraphQl")
	return &model.AuthResponse{AuthToken: &model.AuthToken{AccessToken: response.Token, ExpiredAt: response.Exp}}, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.NewUser) (*models.User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, authorID string, title string, content string) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: CreatePost - createPost"))
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, authorID string, postID *string, parentCommentID *string, content string) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: CreateComment - createComment"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	return users, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	return users[0], nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*models.Post, error) {
	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context) ([]*models.Comment, error) {
	panic(fmt.Errorf("not implemented: Comments - comments"))
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id string) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: Comment - comment"))
}

// Name is the resolver for the name field.
func (r *userResolver) Name(ctx context.Context, obj *models.User) (string, error) {
	panic(fmt.Errorf("not implemented: Name - name"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
