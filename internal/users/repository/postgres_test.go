package repository_test

import (
	"context"
	"database/sql"
	"myHabr/internal/models"
	"myHabr/internal/users/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func (suite *UserRepoTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.Require().NoError(err)
}

func (suite *UserRepoTestSuite) TearDownTest() {
	suite.mock.ExpectClose()
	suite.Require().NoError(suite.db.Close())
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		input      *models.UserSignInUp
		mockExpect func(mock sqlmock.Sqlmock)
		expected   *models.UserCreatedInfo
		expectErr  bool
	}{
		{
			name: "Success",
			input: &models.UserSignInUp{
				Username:     "Auuu",
				PasswordHash: "e38ad214943daad1d64c102faec29de4afe9da3d",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO user_data \(username, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("Auuu", "e38ad214943daad1d64c102faec29de4afe9da3d").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(`SELECT username FROM user_data WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("Auuu"))
			},
			expected:  &models.UserCreatedInfo{ID: 1, Username: "Auuu"},
			expectErr: false,
		},
		{
			name: "Insert Error",
			input: &models.UserSignInUp{
				Username:     "Auuu",
				PasswordHash: "e38ad214943daad1d64c102faec29de4afe9da3d",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO user_data \(username, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("Auuu", "e38ad214943daad1d64c102faec29de4afe9da3d").
					WillReturnError(sql.ErrConnDone)
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Select Error",
			input: &models.UserSignInUp{
				Username:     "Auuu",
				PasswordHash: "e38ad214943daad1d64c102faec29de4afe9da3d",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO user_data \(username, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("Auuu", "e38ad214943daad1d64c102faec29de4afe9da3d").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(`SELECT username FROM user_data WHERE id = \$1`).
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

			repo := repository.NewUserRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			newUser, err := repo.CreateUser(ctx, tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, newUser)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		mockExpect func(mock sqlmock.Sqlmock)
		expected   *models.UserCreatedInfo
		expectedPw string
		expectErr  bool
	}{
		{
			name:     "Success",
			username: "Auuu",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, password_hash FROM user_data WHERE username = \$1`).
					WithArgs("Auuu").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, "e38ad214943daad1d64c102faec29de4afe9da3d"))
			},
			expected:   &models.UserCreatedInfo{ID: 1, Username: "Auuu"},
			expectedPw: "e38ad214943daad1d64c102faec29de4afe9da3d",
			expectErr:  false,
		},
		{
			name:     "User Not Found",
			username: "Auuu",
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, password_hash FROM user_data WHERE username = \$1`).
					WithArgs("Auuu").
					WillReturnError(sql.ErrNoRows)
			},
			expected:   nil,
			expectedPw: "",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewUserRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			user, passwordHash, err := repo.GetUserByLogin(ctx, tt.username)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
				assert.Equal(t, tt.expectedPw, passwordHash)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCheckUser(t *testing.T) {
	tests := []struct {
		name       string
		input      *models.UserSignInUp
		mockExpect func(mock sqlmock.Sqlmock)
		expected   *models.UserCreatedInfo
		expectErr  bool
	}{
		{
			name: "Success",
			input: &models.UserSignInUp{
				Username:     "Auuu",
				PasswordHash: "e38ad214943daad1d64c102faec29de4afe9da3d",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, password_hash FROM user_data WHERE username = \$1`).
					WithArgs("Auuu").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, "e38ad214943daad1d64c102faec29de4afe9da3d"))
			},
			expected:  &models.UserCreatedInfo{ID: 1, Username: "Auuu"},
			expectErr: false,
		},
		{
			name: "Wrong Password",
			input: &models.UserSignInUp{
				Username:     "Auuu",
				PasswordHash: "wrongpassword",
			},
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, password_hash FROM user_data WHERE username = \$1`).
					WithArgs("Auuu").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, "e38ad214943daad1d64c102faec29de4afe9da3d"))
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

			repo := repository.NewUserRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			user, err := repo.CheckUser(ctx, tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUsernameById(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		mockExpect func(mock sqlmock.Sqlmock)
		expected   string
		expectErr  bool
	}{
		{
			name: "Success",
			id:   1,
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT username FROM user_data WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("Auuu"))
			},
			expected:  "Auuu",
			expectErr: false,
		},
		{
			name: "User Not Found",
			id:   1,
			mockExpect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT username FROM user_data WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := repository.NewUserRepo(db)
			ctx := context.Background()

			tt.mockExpect(mock)

			username, err := repo.GetUsernameById(ctx, tt.id)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, username)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
