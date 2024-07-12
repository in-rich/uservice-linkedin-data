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

func TestUpsertUser(t *testing.T) {
	testData := []struct {
		name             string
		publicIdentifier string
		data             *models.UpsertUser

		shouldCreate   bool
		createResponse *entities.User
		createErr      error

		shouldUpdate   bool
		updateResponse *entities.User
		updateErr      error

		shouldCallUpsertProfilePicture bool
		upsertProfilePictureResponse   string
		upsertProfilePictureErr        error

		expect    *models.User
		expectErr error
	}{
		{
			name:             "CreateUser",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName: "first-name-1",
				LastName:  "last-name-1",
			},
			shouldCreate: true,
			createResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
		},
		{
			name:             "UpdateUser",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName: "first-name-2",
				LastName:  "last-name-2",
			},
			shouldCreate: true,
			createErr:    dao.ErrUserAlreadyExists,
			shouldUpdate: true,
			updateResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
		},
		{
			name:             "CreateWithProfilePicture",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName:      "first-name-1",
				LastName:       "last-name-1",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureResponse:   "profile-picture-1",
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
				ProfilePicture:   "profile-picture-1",
			},
		},
		{
			name:             "UpdateWithProfilePicture",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName:      "first-name-2",
				LastName:       "last-name-2",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createErr:    dao.ErrUserAlreadyExists,
			shouldUpdate: true,
			updateResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureResponse:   "profile-picture-1",
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
				ProfilePicture:   "profile-picture-1",
			},
		},
		{
			name:             "CreateWithProfilePictureError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName:      "first-name-1",
				LastName:       "last-name-1",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureErr:        FooErr,
			expectErr:                      FooErr,
		},
		{
			name:             "UpdateWithProfilePictureError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName:      "first-name-2",
				LastName:       "last-name-2",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createErr:    dao.ErrUserAlreadyExists,
			shouldUpdate: true,
			updateResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureErr:        FooErr,
			expectErr:                      FooErr,
		},
		{
			name:             "CreateError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName: "first-name-1",
				LastName:  "last-name-1",
			},
			shouldCreate: true,
			createErr:    FooErr,
			expectErr:    FooErr,
		},
		{
			name:             "UpdateError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName: "first-name-2",
				LastName:  "last-name-2",
			},
			shouldCreate: true,
			createErr:    dao.ErrUserAlreadyExists,
			shouldUpdate: true,
			updateErr:    FooErr,
			expectErr:    FooErr,
		},
		{
			name:             "NoFirstName",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName: "",
				LastName:  "last-name-1",
			},
			expectErr: services.ErrInvalidUpsertUser,
		},
		{
			name:             "NoLastName",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertUser{
				FirstName: "first-name-1",
				LastName:  "",
			},
			expectErr: services.ErrInvalidUpsertUser,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			createUser := daomocks.NewMockCreateUserRepository(t)
			updateUser := daomocks.NewMockUpdateUserRepository(t)
			upsertProfilePicture := daomocks.NewMockUpsertProfilePictureRepository(t)

			if tt.shouldCreate {
				createUser.On("CreateUser", context.TODO(), tt.publicIdentifier, &dao.CreateUserData{
					FirstName: tt.data.FirstName,
					LastName:  tt.data.LastName,
				}).Return(tt.createResponse, tt.createErr)
			}

			if tt.shouldUpdate {
				updateUser.On("UpdateUser", context.TODO(), tt.publicIdentifier, &dao.UpdateUserData{
					FirstName: tt.data.FirstName,
					LastName:  tt.data.LastName,
				}).Return(tt.updateResponse, tt.updateErr)
			}

			if tt.shouldCallUpsertProfilePicture {
				upsertProfilePicture.
					On("UpsertProfilePicture", context.TODO(), tt.publicIdentifier, tt.data.ProfilePicture).
					Return(tt.upsertProfilePictureResponse, tt.upsertProfilePictureErr)
			}

			service := services.NewUpsertUserService(createUser, updateUser, upsertProfilePicture)

			user, err := service.Exec(context.TODO(), tt.publicIdentifier, tt.data)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, user)

			createUser.AssertExpectations(t)
			updateUser.AssertExpectations(t)
			upsertProfilePicture.AssertExpectations(t)
		})
	}
}
