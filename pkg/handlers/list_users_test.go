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

func TestListUsers(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.ListUsersRequest

		listUsersResponse []*models.User
		listUsersErr      error

		expect     *linkedin_data_pb.ListUsersResponse
		expectCode codes.Code
	}{
		{
			name: "ListUsersHandler",
			in: &linkedin_data_pb.ListUsersRequest{
				PublicIdentifiers: []string{"public-identifier-1", "public-identifier-2"},
			},
			listUsersResponse: []*models.User{
				{
					PublicIdentifier: "public-identifier-1",
					FirstName:        "first-name-1",
					LastName:         "last-name-1",
					ProfilePicture:   "profile-picture-1",
				},
				{
					PublicIdentifier: "public-identifier-2",
					FirstName:        "first-name-2",
					LastName:         "last-name-2",
				},
			},
			expect: &linkedin_data_pb.ListUsersResponse{
				Users: []*linkedin_data_pb.User{
					{
						PublicIdentifier:  "public-identifier-1",
						FirstName:         "first-name-1",
						LastName:          "last-name-1",
						ProfilePictureUrl: "profile-picture-1",
					},
					{
						PublicIdentifier: "public-identifier-2",
						FirstName:        "first-name-2",
						LastName:         "last-name-2",
					},
				},
			},
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.ListUsersRequest{
				PublicIdentifiers: []string{"public-identifier-3"},
			},
			listUsersErr: errors.New("internal error"),
			expectCode:   codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockListUsersService(t)
			service.On("Exec", context.TODO(), tt.in.PublicIdentifiers).Return(tt.listUsersResponse, tt.listUsersErr)

			handler := handlers.NewListUsers(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.ListUsers(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
