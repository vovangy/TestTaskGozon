package posts

import (
	"context"
	"myHabr/internal/models"
)

// PostRepo represents the repository interface for posts.
type PostRepo interface {
	CreatePost(ctx context.Context, post *models.PostCreateData) (*models.PostCreateResponse, error)
	CreateComment(ctx context.Context, tx models.Transaction, comment *models.CommentCreateData) (*models.CommentCreateResponse, error)
	BeginTx(ctx context.Context) (models.Transaction, error)
	BlockCommentsOnPost(ctx context.Context, data *models.CommentsBlockRequest) error
	GetCommentsByPostId(ctx context.Context, postId int64) ([]*models.CommentTree, error)
	GetPostById(ctx context.Context, postId int64) (*models.PostCreateResponse, error)
}
