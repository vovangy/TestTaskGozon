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
	postGen "myHabr/internal/posts/delivery/grpc/gen"
	userGen "myHabr/internal/users/delivery/grpc/gen"
	"strconv"
)

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	response, err := r.grpcUserClient.SignUp(ctx, &userGen.SignInUpRequest{Username: input.Username, Password: input.Password})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("User created GraphQl")
	return &model.AuthResponse{AuthToken: &model.AuthToken{AccessToken: response.Token, ExpiredAt: response.Exp}}, nil
}

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	response, err := r.grpcUserClient.Login(ctx, &userGen.SignInUpRequest{Username: input.Username, Password: input.Password})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("User logined GraphQl")
	return &model.AuthResponse{AuthToken: &model.AuthToken{AccessToken: response.Token, ExpiredAt: response.Exp}}, nil
}

// BlockComments is the resolver for the blockComments field.
func (r *mutationResolver) BlockComments(ctx context.Context, postID string) (string, error) {
	id, ok := ctx.Value("userid").(int64)
	if !ok {
		slog.Error("Error with Id")
		return "", fmt.Errorf("Error with id")
	}

	idInt, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	_, err = r.grpcPostClient.BlockCommentsOnPost(ctx, &postGen.BlockCommentsOnPostRequest{UserId: id, PostId: int64(idInt)})
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return "Comments blocked", nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*models.Post, error) {
	id, ok := ctx.Value("userid").(int64)
	if !ok {
		slog.Error("Error with Id")
		return nil, fmt.Errorf("Error with id")
	}

	response, err := r.grpcPostClient.CreatePost(ctx, &postGen.CreatePostRequest{UserId: id, Title: input.Title, Content: input.Content, IsCommented: input.IsCommented})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	responseUser, err := r.grpcUserClient.GetUsernameById(ctx, &userGen.GetUsernameByIdRequest{UserId: id})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Post created GraphQl")
	return &models.Post{ID: strconv.FormatInt(response.PostId, 10), Title: response.Title, Content: response.Content, Author: &models.User{Username: responseUser.Username, ID: strconv.FormatInt(id, 10)}, Comments: nil}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*models.Comment, error) {
	id, ok := ctx.Value("userid").(int64)
	if !ok {
		slog.Error("Error with Id")
		return nil, fmt.Errorf("Error with id")
	}

	postId, err := strconv.Atoi(input.PostID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	parentCommnetId, err := strconv.Atoi(input.ParentCommentID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	response, err := r.grpcPostClient.CreateComment(ctx, &postGen.CreateCommentRequest{UserId: id, PostId: int64(postId), CommentParentId: int64(parentCommnetId), Content: input.Content})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	responseUser, err := r.grpcUserClient.GetUsernameById(ctx, &userGen.GetUsernameByIdRequest{UserId: id})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Comment created GraphQl")
	return &models.Comment{ID: strconv.FormatInt(response.CommentId, 10), Content: response.Content, Author: &models.User{Username: responseUser.Username, ID: strconv.FormatInt(id, 10)}, Replies: nil}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	responseUser, err := r.grpcUserClient.GetUsernameById(ctx, &userGen.GetUsernameByIdRequest{UserId: int64(idInt)})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &models.User{ID: id, Username: responseUser.Username}, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*models.Post, error) {
	responsePosts, err := r.grpcPostClient.GetPosts(ctx, &postGen.GetPostsRequest{})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	posts := []*models.Post{}

	for _, val := range responsePosts.Posts {
		post, err := r.Post(ctx, strconv.FormatInt(val.PostId, 10))
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	responsePost, err := r.grpcPostClient.GetPostById(ctx, &postGen.GetPostByIdRequest{PostId: int64(idInt)})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var comments []*models.Comment

	type StackItem struct {
		Node       *postGen.Comment
		ParentNode *models.Comment
	}

	for i := 0; i < len(responsePost.Comments); i++ {
		stack := []StackItem{{Node: responsePost.Comments[i], ParentNode: nil}}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			responseUser, err := r.grpcUserClient.GetUsernameById(ctx, &userGen.GetUsernameByIdRequest{UserId: current.Node.UserId})
			if err != nil {
				slog.Error(err.Error())
				return nil, err
			}
			convertedComment := &models.Comment{
				ID:      strconv.FormatInt(current.Node.CommentId, 10),
				Content: current.Node.Content,
				Author:  &models.User{ID: strconv.FormatInt(current.Node.UserId, 10), Username: responseUser.Username},
				Replies: []*models.Comment{},
			}

			if current.ParentNode == nil {
				comments = append(comments, convertedComment)
			} else {
				current.ParentNode.Replies = append(current.ParentNode.Replies, convertedComment)
			}

			for i := len(current.Node.Comments) - 1; i >= 0; i-- {
				stack = append(stack, StackItem{Node: current.Node.Comments[i], ParentNode: convertedComment})
			}
		}
	}

	responseUser, err := r.grpcUserClient.GetUsernameById(ctx, &userGen.GetUsernameByIdRequest{UserId: responsePost.UserId})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &models.Post{ID: id, Title: responsePost.Title, Content: responsePost.Content, Author: &models.User{ID: strconv.FormatInt(responsePost.UserId, 10), Username: responseUser.Username}, Comments: comments}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
