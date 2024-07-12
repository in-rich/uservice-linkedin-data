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

func TestUpsertCompany(t *testing.T) {
	testData := []struct {
		name             string
		publicIdentifier string
		data             *models.UpsertCompany

		shouldCreate   bool
		createResponse *entities.Company
		createErr      error

		shouldUpdate   bool
		updateResponse *entities.Company
		updateErr      error

		shouldCallUpsertProfilePicture bool
		upsertProfilePictureResponse   string
		upsertProfilePictureErr        error

		expect    *models.Company
		expectErr error
	}{
		{
			name:             "CreateCompany",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name: "company-1",
			},
			shouldCreate: true,
			createResponse: &entities.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
			},
			expect: &models.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
			},
		},
		{
			name:             "UpdateCompany",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name: "company-2",
			},
			shouldCreate: true,
			createErr:    dao.ErrCompanyAlreadyExists,
			shouldUpdate: true,
			updateResponse: &entities.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-2",
			},
			expect: &models.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-2",
			},
		},
		{
			name:             "CreateWithProfilePicture",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name:           "company-1",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createResponse: &entities.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureResponse:   "profile-picture-1",
			expect: &models.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
				ProfilePicture:   "profile-picture-1",
			},
		},
		{
			name:             "UpdateWithProfilePicture",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name:           "company-2",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createErr:    dao.ErrCompanyAlreadyExists,
			shouldUpdate: true,
			updateResponse: &entities.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-2",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureResponse:   "profile-picture-1",
			expect: &models.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-2",
				ProfilePicture:   "profile-picture-1",
			},
		},
		{
			name:             "CreateWithProfilePictureError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name:           "company-1",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createResponse: &entities.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureErr:        FooErr,
			expectErr:                      FooErr,
		},
		{
			name:             "UpdateWithProfilePictureError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name:           "company-2",
				ProfilePicture: "profile-picture-1",
			},
			shouldCreate: true,
			createErr:    dao.ErrCompanyAlreadyExists,
			shouldUpdate: true,
			updateResponse: &entities.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-2",
			},
			shouldCallUpsertProfilePicture: true,
			upsertProfilePictureErr:        FooErr,
			expectErr:                      FooErr,
		},
		{
			name:             "CreateError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name: "company-1",
			},
			shouldCreate: true,
			createErr:    FooErr,
			expectErr:    FooErr,
		},
		{
			name:             "UpdateError",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name: "company-2",
			},
			shouldCreate: true,
			createErr:    dao.ErrCompanyAlreadyExists,
			shouldUpdate: true,
			updateErr:    FooErr,
			expectErr:    FooErr,
		},
		{
			name:             "NoName",
			publicIdentifier: "public-identifier-1",
			data: &models.UpsertCompany{
				Name: "",
			},
			expectErr: services.ErrInvalidUpsertCompany,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			createCompany := daomocks.NewMockCreateCompanyRepository(t)
			updateCompany := daomocks.NewMockUpdateCompanyRepository(t)
			upsertProfilePicture := daomocks.NewMockUpsertProfilePictureRepository(t)

			if tt.shouldCreate {
				createCompany.On("CreateCompany", context.TODO(), tt.publicIdentifier, &dao.CreateCompanyData{
					Name: tt.data.Name,
				}).Return(tt.createResponse, tt.createErr)
			}

			if tt.shouldUpdate {
				updateCompany.On("UpdateCompany", context.TODO(), tt.publicIdentifier, &dao.UpdateCompanyData{
					Name: tt.data.Name,
				}).Return(tt.updateResponse, tt.updateErr)
			}

			if tt.shouldCallUpsertProfilePicture {
				upsertProfilePicture.
					On("UpsertProfilePicture", context.TODO(), tt.publicIdentifier, tt.data.ProfilePicture).
					Return(tt.upsertProfilePictureResponse, tt.upsertProfilePictureErr)
			}

			service := services.NewUpsertCompanyService(createCompany, updateCompany, upsertProfilePicture)

			user, err := service.Exec(context.TODO(), tt.publicIdentifier, tt.data)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, user)

			createCompany.AssertExpectations(t)
			updateCompany.AssertExpectations(t)
			upsertProfilePicture.AssertExpectations(t)
		})
	}
}
