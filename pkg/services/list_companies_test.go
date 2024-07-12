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

func TestListCompanies(t *testing.T) {
	testData := []struct {
		name string

		publicIdentifiers []string

		listResponse []*entities.Company
		listErr      error

		shouldCallGetProfilePicture bool
		getProfilePictureResponses  map[string]string
		getProfilePictureErrs       map[string]error

		expect    []*models.Company
		expectErr error
	}{
		{
			name: "ListCompanies",
			publicIdentifiers: []string{
				"public-identifier-1",
				"public-identifier-2",
			},
			listResponse: []*entities.Company{
				{
					PublicIdentifier: "public-identifier-1",
					Name:             "company-1",
				},
				{
					PublicIdentifier: "public-identifier-2",
					Name:             "company-2",
				},
			},
			shouldCallGetProfilePicture: true,
			getProfilePictureResponses: map[string]string{
				"public-identifier-1": "profile-picture-1",
				"public-identifier-2": "profile-picture-2",
			},
			expect: []*models.Company{
				{
					PublicIdentifier: "public-identifier-1",
					Name:             "company-1",
					ProfilePicture:   "profile-picture-1",
				},
				{
					PublicIdentifier: "public-identifier-2",
					Name:             "company-2",
					ProfilePicture:   "profile-picture-2",
				},
			},
		},
		{
			name: "MissingProfilePictures",
			publicIdentifiers: []string{
				"public-identifier-1",
				"public-identifier-2",
			},
			listResponse: []*entities.Company{
				{
					PublicIdentifier: "public-identifier-1",
					Name:             "company-1",
				},
				{
					PublicIdentifier: "public-identifier-2",
					Name:             "company-2",
				},
			},
			shouldCallGetProfilePicture: true,
			getProfilePictureResponses: map[string]string{
				"public-identifier-1": "profile-picture-1",
			},
			getProfilePictureErrs: map[string]error{
				"public-identifier-2": dao.ErrProfilePictureNotFound,
			},
			expect: []*models.Company{
				{
					PublicIdentifier: "public-identifier-1",
					Name:             "company-1",
					ProfilePicture:   "profile-picture-1",
				},
				{
					PublicIdentifier: "public-identifier-2",
					Name:             "company-2",
				},
			},
		},
		{
			name: "ListError",
			publicIdentifiers: []string{
				"public-identifier-3",
				"public-identifier-4",
			},
			listErr:   FooErr,
			expectErr: FooErr,
		},
		{
			name: "ProfilePictureError",
			publicIdentifiers: []string{
				"public-identifier-1",
				"public-identifier-2",
			},
			listResponse: []*entities.Company{
				{
					PublicIdentifier: "public-identifier-1",
					Name:             "company-1",
				},
				{
					PublicIdentifier: "public-identifier-2",
					Name:             "company-2",
				},
			},
			shouldCallGetProfilePicture: true,
			getProfilePictureResponses: map[string]string{
				"public-identifier-1": "profile-picture-1",
			},
			getProfilePictureErrs: map[string]error{
				"public-identifier-2": FooErr,
			},
			expectErr: FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			listCompanies := daomocks.NewMockListCompaniesRepository(t)
			listCompanies.On("ListCompanies", context.TODO(), tt.publicIdentifiers).Return(tt.listResponse, tt.listErr)

			getProfilePicture := daomocks.NewMockGetProfilePictureRepository(t)
			if tt.shouldCallGetProfilePicture {
				for _, user := range tt.listResponse {
					getProfilePicture.
						On("GetProfilePicture", context.TODO(), user.PublicIdentifier).
						Return(tt.getProfilePictureResponses[user.PublicIdentifier], tt.getProfilePictureErrs[user.PublicIdentifier])
				}
			}

			service := services.NewListCompaniesService(listCompanies, getProfilePicture)
			res, err := service.Exec(context.TODO(), tt.publicIdentifiers)

			require.Equal(t, tt.expect, res)
			require.Equal(t, tt.expectErr, err)

			listCompanies.AssertExpectations(t)
			getProfilePicture.AssertExpectations(t)
		})
	}
}
