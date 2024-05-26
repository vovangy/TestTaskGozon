package grpc

import (
	"context"
	"log/slog"
	"myHabr/internal/models"
	"myHabr/internal/posts"
	genPosts "myHabr/internal/posts/delivery/grpc/gen"
)

type PostsServerHandler struct {
	genPosts.PostServer
	uc posts.PostUsecase
}

func NewPostsServerHandler(uc posts.PostUsecase) *PostsServerHandler {
	return &PostsServerHandler{uc: uc}
}

func (h *PostsServerHandler) CreatePost(ctx context.Context, req *genPosts.CreatePostRequest) (*genPosts.CreatePostResponse, error) {

	post, err := h.uc.CreatePost(ctx, &models.PostCreateData{UserId: req.UserId, IsCommented: req.IsCommented, Title: req.Content, Content: req.Content})

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Success created post Grpc")
	return &genPosts.CreatePostResponse{UserId: post.UserId, Title: post.Title, Content: post.Content, PostId: post.ID, IsCommented: post.IsCommented}, nil
}

func (h *PostsServerHandler) CreateComment(ctx context.Context, req *genPosts.CreateCommentRequest) (*genPosts.CreateCommentResponse, error) {

	comment, err := h.uc.CreateComment(ctx, &models.CommentCreateData{UserId: req.UserId, PostId: req.PostId, CommentParentId: req.CommentParentId, Content: req.Content})

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Success created comment Grpc")
	return &genPosts.CreateCommentResponse{CommentId: comment.CommentId, UserId: comment.UserId, PostId: comment.PostId, CommentParentId: comment.CommentParentId, Content: comment.Content}, nil
}

func (h *PostsServerHandler) BlockCommentsOnPost(ctx context.Context, req *genPosts.BlockCommentsOnPostRequest) (*genPosts.BlockCommentsOnPostResponse, error) {

	err := h.uc.BlockCommentsOnPost(ctx, &models.CommentsBlockRequest{UserId: req.UserId, PostId: req.PostId})

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Success blocked comments on post Grpc")
	return &genPosts.BlockCommentsOnPostResponse{}, nil
}

func (h *PostsServerHandler) GetPostById(ctx context.Context, req *genPosts.GetPostByIdRequest) (*genPosts.GetPostByIdResponse, error) {

	post, err := h.uc.GetPostById(ctx, req.PostId)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var comments []*genPosts.Comment

	type StackItem struct {
		Node       *models.CommentTree
		ParentNode *genPosts.Comment
	}

	for i := 0; i < len(post.Comments); i++ {
		stack := []StackItem{{Node: post.Comments[i], ParentNode: nil}}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			convertedComment := &genPosts.Comment{
				CommentId: current.Node.Comment.CommentId,
				UserId:    current.Node.Comment.UserId,
				PostId:    current.Node.Comment.PostId,
				Content:   current.Node.Comment.Content,
				Comments:  []*genPosts.Comment{},
			}

			if current.ParentNode == nil {
				comments = append(comments, convertedComment)
			} else {
				current.ParentNode.Comments = append(current.ParentNode.Comments, convertedComment)
			}

			for i := len(current.Node.Replies) - 1; i >= 0; i-- {
				stack = append(stack, StackItem{Node: current.Node.Replies[i], ParentNode: convertedComment})
			}
		}
	}

	slog.Info("Success getting post Grpc")
	return &genPosts.GetPostByIdResponse{UserId: post.Post.UserId, Title: post.Post.Title, Content: post.Post.Content, PostId: post.Post.ID, IsCommented: post.Post.IsCommented, Comments: comments}, nil
}
