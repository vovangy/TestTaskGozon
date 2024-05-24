package usecase

import (
	"context"
	"log/slog"
	"myHabr/internal/middleware/jwt"
	"myHabr/internal/models"
	"myHabr/internal/pkg"
	"myHabr/internal/users"
	"time"
)

// UserUsecase represents the usecase for users.
type UserUsecase struct {
	repo users.UserRepo
}

// NewAuthUsecase creates a new instance of AuthUsecase.
func NewUserUsecase(repo users.UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// SignUp handles the user registration process.
func (u *UserUsecase) SignUp(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, string, time.Time, error) {
	data.PasswordHash = pkg.GenerateHashString(data.PasswordHash)
	userResponse, err := u.repo.CreateUser(ctx, data)

	if err != nil {
		slog.Error(err.Error())
		return nil, "", time.Now(), err
	}

	token, exp, err := jwt.GenerateToken(userResponse)
	if err != nil {
		slog.Error(err.Error())
		return nil, "", time.Now(), err
	}

	slog.Info("Created user Usecase")
	return userResponse, token, exp, nil
}

// Login handles the user login process.
func (u *UserUsecase) Login(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, string, time.Time, error) {
	data.PasswordHash = pkg.GenerateHashString(data.PasswordHash)
	user, err := u.repo.CheckUser(ctx, data)
	if err != nil {
		slog.Error(err.Error())
		return nil, "", time.Now(), err
	}

	token, exp, err := jwt.GenerateToken(user)
	if err != nil {
		slog.Error(err.Error())
		return nil, "", time.Now(), err
	}

	slog.Info("Success login Usecase")
	return user, token, exp, nil
}

// GetUsernameById handles the user get username by id process.
func (u *UserUsecase) GetUsernameById(ctx context.Context, id int64) (string, error) {
	username, err := u.repo.GetUsernameById(ctx, id)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("Success get username by id Usecase")
	return username, nil
}
