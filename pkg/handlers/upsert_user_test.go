package handlers_test

import (
	"context"
	"errors"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	servicesmocks "github.com/in-rich/uservice-linkedin-data/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestUpsertUser(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.UpsertUserRequest

		upsertResponse *models.User
		upsertErr      error

		expect     *linkedin_data_pb.User
		expectCode codes.Code
	}{
		{
			name: "UpsertUser",
			in: &linkedin_data_pb.UpsertUserRequest{
				PublicIdentifier:     "public-identifier-1",
				FirstName:            "first-name-1",
				LastName:             "last-name-1",
				ProfilePictureBase64: "profile-picture-1-content",
			},
			upsertResponse: &models.User{
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
		},
		{
			name: "InvalidUpsertUser",
			in: &linkedin_data_pb.UpsertUserRequest{
				PublicIdentifier: "public-identifier-2",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
			upsertErr:  services.ErrInvalidUpsertUser,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.UpsertUserRequest{
				PublicIdentifier: "public-identifier-3",
				FirstName:        "first-name-3",
				LastName:         "last-name-3",
			},
			upsertErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpsertUserService(t)
			service.
				On("Exec", context.TODO(), tt.in.PublicIdentifier, &models.UpsertUser{
					FirstName:      tt.in.FirstName,
					LastName:       tt.in.LastName,
					ProfilePicture: tt.in.ProfilePictureBase64,
				}).
				Return(tt.upsertResponse, tt.upsertErr)

			handler := handlers.NewUpsertUser(service)
			resp, err := handler.UpsertUser(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
