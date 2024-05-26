package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"myHabr/internal/models"
)

// PostRepo represents a repository for Posts.
type PostRepo struct {
	db *sql.DB
}

// NewPostRepo creates a new instance of PostRepo.
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) BeginTx(ctx context.Context) (models.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Transaction Started")
	return tx, nil
}

// CreatePost creates a new post in the database.
func (r *PostRepo) CreatePost(ctx context.Context, post *models.PostCreateData) (*models.PostCreateResponse, error) {

	insert := `INSERT INTO post (user_id, is_commented, title, content) VALUES ($1, $2, $3, $4) RETURNING id`
	var lastInsertID int64

	if err := r.db.QueryRowContext(ctx, insert, post.UserId, post.IsCommented, post.Title, post.Content).Scan(&lastInsertID); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	query := `SELECT user_id, is_commented, title, content FROM post WHERE id = $1`

	res := r.db.QueryRow(query, lastInsertID)

	newPost := &models.PostCreateResponse{ID: lastInsertID}
	if err := res.Scan(&newPost.UserId, &newPost.IsCommented, &newPost.Title, &newPost.Content); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Created post Repository")
	return newPost, nil
}

// CreateComment creates a new comment in the database.
func (r *PostRepo) CreateComment(ctx context.Context, tx models.Transaction, comment *models.CommentCreateData) (*models.CommentCreateResponse, error) {
	insertComment := `INSERT INTO comment (user_id, post_id, content) VALUES ($1, $2, $3) RETURNING id`
	var lastInsertID int64

	if err := tx.QueryRowContext(ctx, insertComment, comment.UserId, comment.PostId, comment.Content).Scan(&lastInsertID); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	insertCommentClosure := `INSERT INTO comment_closure (descendant_id) VALUES ($1)`
	ids := []interface{}{lastInsertID}
	if comment.CommentParentId != 0 {
		insertCommentClosure = `INSERT INTO comment_closure (descendant_id, ancestor_id) VALUES ($1, $2)`
		ids = append(ids, comment.CommentParentId)
	}

	if _, err := tx.ExecContext(ctx, insertCommentClosure, ids...); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	query := `SELECT c.user_id, c.post_id, c.content, COALESCE(cc.ancestor_id, 0) FROM comment as c LEFT JOIN comment_closure AS cc ON c.id=cc.descendant_id WHERE c.id = $1`

	res := tx.QueryRowContext(ctx, query, lastInsertID)

	newComment := &models.CommentCreateResponse{CommentId: lastInsertID}
	if err := res.Scan(&newComment.UserId, &newComment.PostId, &newComment.Content, &newComment.CommentParentId); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Created comment Repository")

	return newComment, nil
}

// BlockCommentsOnPost block a new comments on post in the database.
func (r *PostRepo) BlockCommentsOnPost(ctx context.Context, data *models.CommentsBlockRequest) error {
	updatePost := `UPDATE post SET is_commented = false WHERE user_id = $1 AND id = $2`

	if _, err := r.db.ExecContext(ctx, updatePost, data.UserId, data.PostId); err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Comments succesfully blocked Repository")

	return nil
}

// GetAllPostIds возвращает список всех ID постов из базы данных.
func (r *PostRepo) GetAllPostIds(ctx context.Context) ([]int64, error) {
	var postIds []int64

	query := `SELECT id FROM post`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postId int64
		if err := rows.Scan(&postId); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		postIds = append(postIds, postId)
	}
	if err := rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("All Post IDs successfully retrieved from Repository")

	return postIds, nil
}

// GetPostById getting a post from the database.
func (r *PostRepo) GetPostById(ctx context.Context, postId int64) (*models.PostCreateResponse, error) {
	selectPost := `SELECT user_id, is_commented, title, content FROM post WHERE id = $1`

	res := r.db.QueryRowContext(ctx, selectPost, postId)

	post := &models.PostCreateResponse{ID: postId}
	if err := res.Scan(&post.UserId, &post.IsCommented, &post.Title, &post.Content); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Post succesfully got Repository")

	return post, nil
}

// getDirectDescendants getting a descendants comments from the database.
func (r *PostRepo) GetDirectDescendants(commentID int64) ([]models.CommentCreateResponse, error) {
	rows, err := r.db.Query(`
        SELECT c.id, c.content, c.user_id, c.post_id
        FROM comment c
        JOIN comment_closure cc ON c.id = cc.descendant_id
        WHERE cc.ancestor_id = $1
    `, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentCreateResponse
	for rows.Next() {
		comment := models.CommentCreateResponse{CommentParentId: commentID}
		if err := rows.Scan(&comment.CommentId, &comment.Content, &comment.UserId, &comment.PostId); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *PostRepo) BuildCommentTree(node *models.CommentTree) error {
	descendants, err := r.GetDirectDescendants(node.Comment.CommentId)
	if err != nil {
		return err
	}

	for _, desc := range descendants {
		childNode := &models.CommentTree{
			Comment: desc,
		}
		node.Replies = append(node.Replies, childNode)

		err = r.BuildCommentTree(node.Replies[len(node.Replies)-1])
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCommentsByPostId getting a comments on post from the database by post id.
func (r *PostRepo) GetCommentsByPostId(ctx context.Context, postId int64) ([]*models.CommentTree, error) {
	selectComments := `SELECT c.id, c.content, c.user_id FROM comment AS c JOIN comment_closure AS cc ON cc.ancestor_id IS NULL AND cc.descendant_id=c.id WHERE c.post_id = $1`

	rows, err := r.db.QueryContext(ctx, selectComments, postId)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	defer rows.Close()

	comments := []*models.CommentTree{}

	for rows.Next() {
		comment := &models.CommentCreateResponse{PostId: postId, CommentParentId: 0}
		err := rows.Scan(&comment.CommentId, &comment.Content, &comment.UserId)

		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		commentsTree := &models.CommentTree{Comment: *comment, Replies: []*models.CommentTree{}}

		if err = r.BuildCommentTree(commentsTree); err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		comments = append(comments, commentsTree)
	}

	if err := rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Comments succesfully got Repository")

	return comments, nil
}
