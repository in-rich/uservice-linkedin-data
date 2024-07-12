package services_test

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	daomocks "github.com/in-rich/uservice-linkedin-data/pkg/dao/mocks"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUser(t *testing.T) {
	testData := []struct {
		name string

		publicIdentifier string

		getResponse *entities.User
		getErr      error

		shouldCallGetProfilePicture bool
		getProfilePictureResponse   string
		getProfilePictureErr        error

		expect    *models.User
		expectErr error
	}{
		{
			name:             "GetUser",
			publicIdentifier: "public-identifier-1",
			getResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
			shouldCallGetProfilePicture: true,
			getProfilePictureResponse:   "profile-picture-1",
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
				ProfilePicture:   "profile-picture-1",
			},
		},
		{
			name:             "NoProfilePicture",
			publicIdentifier: "public-identifier-1",
			getResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
			shouldCallGetProfilePicture: true,
			getProfilePictureErr:        dao.ErrProfilePictureNotFound,
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
		},
		{
			name:             "GetError",
			publicIdentifier: "public-identifier-1",
			getErr:           FooErr,
			expectErr:        FooErr,
		},
		{
			name:             "GetProfilePictureError",
			publicIdentifier: "public-identifier-1",
			getResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
			shouldCallGetProfilePicture: true,
			getProfilePictureErr:        FooErr,
			expectErr:                   FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			getUser := daomocks.NewMockGetUserRepository(t)
			getUser.On("GetUser", context.TODO(), tt.publicIdentifier).Return(tt.getResponse, tt.getErr)

			getProfilePicture := daomocks.NewMockGetProfilePictureRepository(t)
			if tt.shouldCallGetProfilePicture {
				getProfilePicture.
					On("GetProfilePicture", context.TODO(), tt.publicIdentifier).
					Return(tt.getProfilePictureResponse, tt.getProfilePictureErr)
			}

			service := services.NewGetUserService(getUser, getProfilePicture)
			res, err := service.Exec(context.TODO(), tt.publicIdentifier)

			require.Equal(t, tt.expect, res)
			require.Equal(t, tt.expectErr, err)

			getUser.AssertExpectations(t)
			getProfilePicture.AssertExpectations(t)
		})
	}
}
