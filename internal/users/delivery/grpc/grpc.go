package grpc

import (
	"context"
	"errors"
	"log/slog"
	"myHabr/internal/models"
	"myHabr/internal/users"
	genUsers "myHabr/internal/users/delivery/grpc/gen"
)

type UsersServerHandler struct {
	genUsers.UserServer
	uc users.UserUsecase
}

func NewUsersServerHandler(uc users.UserUsecase) *UsersServerHandler {
	return &UsersServerHandler{uc: uc}
}

func (h *UsersServerHandler) SignUp(ctx context.Context, req *genUsers.SignInUpRequest) (*genUsers.SignUpInResponse, error) {

	_, token, exp, err := h.uc.SignUp(ctx, &models.UserSignInUp{Username: req.Username, PasswordHash: req.Password})

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	layout := "2006-01-02 15:04:05"
	dateString := exp.Format(layout)

	slog.Info("Success sign up user Grpc")
	return &genUsers.SignUpInResponse{Token: token, Exp: dateString}, nil
}

func (h *UsersServerHandler) Login(ctx context.Context, req *genUsers.SignInUpRequest) (*genUsers.SignUpInResponse, error) {
	_, token, exp, err := h.uc.Login(ctx, &models.UserSignInUp{Username: req.Username, PasswordHash: req.Password})

	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("error login")
	}

	layout := "2006-01-02 15:04:05"
	dateString := exp.Format(layout)

	slog.Info("Success login user Grpc")
	return &genUsers.SignUpInResponse{Token: token, Exp: dateString}, nil
}

func (h *UsersServerHandler) GetUsernameById(ctx context.Context, req *genUsers.GetUsernameByIdRequest) (*genUsers.GetUsernameByIdResponse, error) {
	username, err := h.uc.GetUsernameById(ctx, req.UserId)

	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("Error getting username")
	}

	slog.Info("Success login user Grpc")
	return &genUsers.GetUsernameByIdResponse{Username: username}, nil
}
