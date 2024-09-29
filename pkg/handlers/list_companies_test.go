package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	servicesmocks "github.com/in-rich/uservice-linkedin-data/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestListCompanies(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.ListCompaniesRequest

		listCompaniesResponse []*models.Company
		listCompaniesErr      error

		expect     *linkedin_data_pb.ListCompaniesResponse
		expectCode codes.Code
	}{
		{
			name: "ListCompaniesHandler",
			in: &linkedin_data_pb.ListCompaniesRequest{
				PublicIdentifiers: []string{"public-identifier-1", "public-identifier-2"},
			},
			listCompaniesResponse: []*models.Company{
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
			expect: &linkedin_data_pb.ListCompaniesResponse{
				Companies: []*linkedin_data_pb.Company{
					{
						PublicIdentifier:  "public-identifier-1",
						Name:              "company-1",
						ProfilePictureUrl: "profile-picture-1",
					},
					{
						PublicIdentifier: "public-identifier-2",
						Name:             "company-2",
					},
				},
			},
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.ListCompaniesRequest{
				PublicIdentifiers: []string{"public-identifier-3"},
			},
			listCompaniesErr: errors.New("internal error"),
			expectCode:       codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockListCompaniesService(t)
			service.On("Exec", context.TODO(), tt.in.PublicIdentifiers).Return(tt.listCompaniesResponse, tt.listCompaniesErr)

			handler := handlers.NewListCompanies(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.ListCompanies(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
