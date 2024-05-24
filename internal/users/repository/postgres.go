package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"myHabr/internal/models"
)

// UserRepo represents a repository for Users.
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new instance of AuthRepo.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateUser creates a new user in the database.
func (r *UserRepo) CreateUser(ctx context.Context, user *models.UserSignInUp) (*models.UserCreatedInfo, error) {

	insert := `INSERT INTO user_data (username, password_hash) VALUES ($1, $2) RETURNING id`
	var lastInsertID int64

	if err := r.db.QueryRowContext(ctx, insert, user.Username, user.PasswordHash).Scan(&lastInsertID); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	query := `SELECT username FROM user_data WHERE id = $1`

	res := r.db.QueryRow(query, lastInsertID)

	newUser := &models.UserCreatedInfo{ID: lastInsertID}
	if err := res.Scan(&newUser.Username); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Created user Repository")
	return newUser, nil
}

// GetUserByLogin retrieves a user from the database by their login.
func (r *UserRepo) GetUserByLogin(ctx context.Context, username string) (*models.UserCreatedInfo, string, error) {
	query := `SELECT id, password_hash FROM user_data WHERE username = $1`

	res := r.db.QueryRowContext(ctx, query, username)

	var passwordHash string
	user := &models.UserCreatedInfo{Username: username}
	if err := res.Scan(&user.ID, &passwordHash); err != nil {
		slog.Error(err.Error())
		return nil, "", err
	}

	slog.Info("Get User By login succes Repository")
	return user, passwordHash, nil
}

// CheckUser checks if the user with the given login and password hash exists in the database.
func (r *UserRepo) CheckUser(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, error) {
	user, passwordHash, err := r.GetUserByLogin(ctx, data.Username)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if passwordHash != data.PasswordHash {
		return nil, errors.New("wrong password")
	}

	slog.Info("success checkUser Repository")
	return user, nil
}

func (r *UserRepo) GetUsernameById(ctx context.Context, id int64) (string, error) {
	query := `SELECT username FROM user_data WHERE id = $1`

	res := r.db.QueryRow(query, id)

	username := ""
	if err := res.Scan(&username); err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("Success get username by user id Repository")
	return username, nil
}
