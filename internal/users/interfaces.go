package users

import (
	"context"
	"myHabr/internal/models"
	"time"
)

const (
	SignUpMethod           = "SignUp"
	LoginMethod            = "Login"
	CheckAuthMethod        = "CheckAuth"
	GetUserLevelByIdMethod = "GetUserLevelById"
	CreateUserMethod       = "CreateUser"
	CheckUserMethod        = "CheckUser"
	GetUserByLoginMethod   = "GetUserByLogin"
	BeginTxMethod          = "BeginTx"
)

// UserUsecase represents the usecase interface for authentication.
type UserUsecase interface {
	SignUp(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, string, time.Time, error)
	Login(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, string, time.Time, error)
	GetUsernameById(ctx context.Context, id int64) (string, error)
}

// UserRepo represents the repository interface for users.
type UserRepo interface {
	CreateUser(ctx context.Context, user *models.UserSignInUp) (*models.UserCreatedInfo, error)
	GetUserByLogin(ctx context.Context, username string) (*models.UserCreatedInfo, string, error)
	CheckUser(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, error)
	GetUsernameById(ctx context.Context, id int64) (string, error)
}
