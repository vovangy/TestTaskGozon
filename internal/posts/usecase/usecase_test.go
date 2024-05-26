package usecase_test

import (
	"context"
	"errors"
	"myHabr/internal/models"
	mock_posts "myHabr/internal/posts/mock"
	"myHabr/internal/posts/usecase"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_posts.NewMockPostRepo(ctrl)

	uc := usecase.NewPostUsecase(mockRepo)

	testCases := []struct {
		name             string
		data             *models.PostCreateData
		repoExpectation  func(repo *mock_posts.MockPostRepo, ctx context.Context, data *models.PostCreateData)
		expectedResponse *models.PostCreateResponse
		expectedError    error
	}{
		{
			name: "Success",
			data: &models.PostCreateData{
				Title:   "Test Post",
				Content: "Test content for the post",
			},
			repoExpectation: func(repo *mock_posts.MockPostRepo, ctx context.Context, data *models.PostCreateData) {
				resp := &models.PostCreateResponse{
					ID:      1,
					Title:   data.Title,
					Content: data.Content,
				}
				repo.EXPECT().CreatePost(ctx, data).Return(resp, nil)
			},
			expectedError:    nil,
			expectedResponse: &models.PostCreateResponse{ID: 1, Title: "Test Post", Content: "Test content for the post"},
		},
		{
			name: "RepositoryError",
			data: &models.PostCreateData{
				Title:   "Test Post",
				Content: "Test content for the post",
			},
			repoExpectation: func(repo *mock_posts.MockPostRepo, ctx context.Context, data *models.PostCreateData) {
				repo.EXPECT().CreatePost(ctx, data).Return(nil, errors.New("repository error"))
			},
			expectedError:    errors.New("repository error"),
			expectedResponse: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.repoExpectation(mockRepo, context.Background(), tc.data)

			resp, err := uc.CreatePost(context.Background(), tc.data)

			if (err != nil || tc.expectedError != nil) && (err == nil || tc.expectedError == nil || err.Error() != tc.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(resp, tc.expectedResponse) {
				t.Errorf("Expected response: %v, got: %v", tc.expectedResponse, resp)
			}
		})
	}
}
