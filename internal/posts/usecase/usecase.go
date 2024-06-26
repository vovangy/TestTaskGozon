package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"myHabr/internal/models"
	"myHabr/internal/posts"
)

// PostUsecase represents the usecase for post using.
type PostUsecase struct {
	repo posts.PostRepo
}

// NewPostUsecase creates a new instance of PostUsecase.
func NewPostUsecase(repo posts.PostRepo) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (u *PostUsecase) CreatePost(ctx context.Context, data *models.PostCreateData) (*models.PostCreateResponse, error) {

	post, err := u.repo.CreatePost(ctx, data)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return post, nil
}

func (u *PostUsecase) CreateComment(ctx context.Context, data *models.CommentCreateData) (*models.CommentCreateResponse, error) {
	post, err := u.repo.GetPostById(ctx, data.PostId)
	if post.IsCommented != true {
		return nil, fmt.Errorf("Post is not commented")
	}

	tx, err := u.repo.BeginTx(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			slog.Error(err.Error())
		}
	}()

	comment, err := u.repo.CreateComment(ctx, tx, data)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Transaction Succesfully Commited")
	return comment, nil
}

func (u *PostUsecase) BlockCommentsOnPost(ctx context.Context, data *models.CommentsBlockRequest) error {
	err := u.repo.BlockCommentsOnPost(ctx, data)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Comments Succesfully blocked Usecase")
	return nil
}

func (u *PostUsecase) GetPosts(ctx context.Context) ([]*models.PostResponse, error) {
	postIds, err := u.repo.GetAllPostIds(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	posts := []*models.PostResponse{}

	for _, val := range postIds {
		post, err := u.GetPostById(ctx, val)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}

	slog.Info("Comments Succesfully blocked Usecase")
	return posts, nil
}

func (u *PostUsecase) GetPostById(ctx context.Context, postId int64) (*models.PostResponse, error) {
	post, err := u.repo.GetPostById(ctx, postId)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	comments, err := u.repo.GetCommentsByPostId(ctx, postId)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Comments Succesfully got Usecase")
	return &models.PostResponse{Post: *post, Comments: comments}, nil
}
