package repository_test

import (
	"context"
	"database/sql"
	"myHabr/internal/models"
	"myHabr/internal/posts/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name       string
		input      *models.PostCreateData
		mockExpect func(mock sqlmock.Sqlmock)
		expected   *models.PostCreateResponse
		expectErr  bool
	}{
		{
			name: "Success",
			input: &models.PostCreateData{
				UserId:      1,
				IsCommented: true,
				Title:       "Test Post",
				Content:     "This is a test post",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO post \(user_id, is_commented, title, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs(1, true, "Test Post", "This is a test post").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(`SELECT user_id, is_commented, title, content FROM post WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "is_commented", "title", "content"}).AddRow(1, true, "Test Post", "This is a test post"))
			},
			expected: &models.PostCreateResponse{
				ID:          1,
				UserId:      1,
				IsCommented: true,
				Title:       "Test Post",
				Content:     "This is a test post",
			},
			expectErr: false,
		},
		{
			name: "Insert Error",
			input: &models.PostCreateData{
				UserId:      1,
				IsCommented: true,
				Title:       "Test Post",
				Content:     "This is a test post",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO post \(user_id, is_commented, title, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs(1, true, "Test Post", "This is a test post").
					WillReturnError(sql.ErrConnDone)
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Select Error",
			input: &models.PostCreateData{
				UserId:      1,
				IsCommented: true,
				Title:       "Test Post",
				Content:     "This is a test post",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO post \(user_id, is_commented, title, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs(1, true, "Test Post", "This is a test post").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(`SELECT user_id, is_commented, title, content FROM post WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewPostRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			newPost, err := repo.CreatePost(ctx, tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, newPost)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateComment(t *testing.T) {
	tests := []struct {
		name       string
		input      *models.CommentCreateData
		mockExpect func(mock sqlmock.Sqlmock)
		expected   *models.CommentCreateResponse
		expectErr  bool
	}{
		{
			name: "Success",
			input: &models.CommentCreateData{
				UserId:          1,
				PostId:          1,
				Content:         "Test Comment",
				CommentParentId: 0,
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO comment \(user_id, post_id, content\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
					WithArgs(1, 1, "Test Comment").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec(`INSERT INTO comment_closure \(descendant_id\) VALUES \(\$1\)`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`SELECT c.user_id, c.post_id, c.content, COALESCE\(cc.ancestor_id, 0\) FROM comment as c LEFT JOIN comment_closure AS cc ON c.id=cc.descendant_id WHERE c.id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "post_id", "content", "ancestor_id"}).AddRow(1, 1, "Test Comment", 0))
				mock.ExpectCommit()
			},
			expected: &models.CommentCreateResponse{
				CommentId:       1,
				UserId:          1,
				PostId:          1,
				Content:         "Test Comment",
				CommentParentId: 0,
			},
			expectErr: false,
		},
		{
			name: "Insert Comment Error",
			input: &models.CommentCreateData{
				UserId:          1,
				PostId:          1,
				Content:         "Test Comment",
				CommentParentId: 0,
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO comment \(user_id, post_id, content\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
					WithArgs(1, 1, "Test Comment").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewPostRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			assert.NotNil(t, tx)

			newComment, err := repo.CreateComment(ctx, tx, tt.input)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, newComment)
				tx.Rollback()
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, newComment)
				tx.Commit()
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %s", err)
			}
		})
	}
}

func TestBlockCommentsOnPost(t *testing.T) {
	tests := []struct {
		name       string
		input      *models.CommentsBlockRequest
		mockExpect func(mock sqlmock.Sqlmock)
		expectErr  bool
	}{
		{
			name: "Success",
			input: &models.CommentsBlockRequest{
				UserId: 1,
				PostId: 1,
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE post SET is_commented = false WHERE user_id = \$1 AND id = \$2`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name: "Update Error",
			input: &models.CommentsBlockRequest{
				UserId: 1,
				PostId: 1,
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE post SET is_commented = false WHERE user_id = \$1 AND id = \$2`).
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewPostRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			err = repo.BlockCommentsOnPost(ctx, tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllPostIds(t *testing.T) {
	tests := []struct {
		name       string
		mockExpect func(mock sqlmock.Sqlmock)
		expected   []int64
		expectErr  bool
	}{
		{
			name: "Success",
			mockExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2).
					AddRow(3)
				mock.ExpectQuery(`SELECT id FROM post`).
					WillReturnRows(rows)
			},
			expected:  []int64{1, 2, 3},
			expectErr: false,
		},
		{
			name: "Query Error",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id FROM post`).
					WillReturnError(sql.ErrConnDone)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewPostRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			postIds, err := repo.GetAllPostIds(ctx)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, postIds)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, postIds)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPostById(t *testing.T) {
	tests := []struct {
		name       string
		postId     int64
		mockExpect func(mock sqlmock.Sqlmock)
		expected   *models.PostCreateResponse
		expectErr  bool
	}{
		{
			name:   "Success",
			postId: 1,
			mockExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id", "is_commented", "title", "content"}).
					AddRow(1, true, "Title 1", "Content 1")
				mock.ExpectQuery(`SELECT user_id, is_commented, title, content FROM post WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: &models.PostCreateResponse{
				ID:          1,
				UserId:      1,
				IsCommented: true,
				Title:       "Title 1",
				Content:     "Content 1",
			},
			expectErr: false,
		},
		{
			name:   "Query Error",
			postId: 1,
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT user_id, is_commented, title, content FROM post WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewPostRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			post, err := repo.GetPostById(ctx, tt.postId)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, post)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetDirectDescendants(t *testing.T) {
	tests := []struct {
		name       string
		commentID  int64
		mockExpect func(mock sqlmock.Sqlmock)
		expected   []models.CommentCreateResponse
		expectErr  bool
	}{
		{
			name:      "Success",
			commentID: 1,
			mockExpect: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "content", "user_id", "post_id"}).
					AddRow(1, "Content 1", 1, 1).
					AddRow(2, "Content 2", 2, 1)
				mock.ExpectQuery(`
					SELECT c.id, c.content, c.user_id, c.post_id
					FROM comment c
					JOIN comment_closure cc ON c.id = cc.descendant_id
					WHERE cc.ancestor_id = \$1
				`).WithArgs(1).WillReturnRows(rows)
			},
			expected: []models.CommentCreateResponse{
				{CommentId: 1, Content: "Content 1", UserId: 1, PostId: 1, CommentParentId: 1},
				{CommentId: 2, Content: "Content 2", UserId: 2, PostId: 1, CommentParentId: 1},
			},
			expectErr: false,
		},
		{
			name:      "Query Error",
			commentID: 1,
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT c.id, c.content, c.user_id, c.post_id
					FROM comment c
					JOIN comment_closure cc ON c.id = cc.descendant_id
					WHERE cc.ancestor_id = \$1
				`).WithArgs(1).WillReturnError(sql.ErrConnDone)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewPostRepo(db)

			tt.mockExpect(mock)

			comments, err := repo.GetDirectDescendants(tt.commentID)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, comments)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, comments)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
