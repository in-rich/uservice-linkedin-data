package handlers_test

import (
	"context"
	"errors"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	servicesmocks "github.com/in-rich/uservice-linkedin-data/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestGetUserData(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.GetUserRequest

		getResponse *models.User
		getErr      error

		expect     *linkedin_data_pb.User
		expectCode codes.Code
	}{
		{
			name: "GetUserHandler",
			in: &linkedin_data_pb.GetUserRequest{
				PublicIdentifier: "public-identifier-1",
			},
			getResponse: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-1",
				LastName:         "last-name-1",
				ProfilePicture:   "profile-picture-1",
			},
			expect: &linkedin_data_pb.User{
				PublicIdentifier:  "public-identifier-1",
				FirstName:         "first-name-1",
				LastName:          "last-name-1",
				ProfilePictureUrl: "profile-picture-1",
			},
			expectCode: codes.OK,
		},
		{
			name: "UserNotFound",
			in: &linkedin_data_pb.GetUserRequest{
				PublicIdentifier: "public-identifier-2",
			},
			getErr:     dao.ErrUserNotFound,
			expectCode: codes.NotFound,
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.GetUserRequest{
				PublicIdentifier: "public-identifier-3",
			},
			getErr:     errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetUserService(t)
			service.On("Exec", context.TODO(), tt.in.PublicIdentifier).Return(tt.getResponse, tt.getErr)

			handler := handlers.NewGetUser(service)
			resp, err := handler.GetUser(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
