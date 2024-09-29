package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	servicesmocks "github.com/in-rich/uservice-linkedin-data/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestUpsertCompany(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.UpsertCompanyRequest

		upsertResponse *models.Company
		upsertErr      error

		expect     *linkedin_data_pb.Company
		expectCode codes.Code
	}{
		{
			name: "UpsertCompany",
			in: &linkedin_data_pb.UpsertCompanyRequest{
				PublicIdentifier:     "public-identifier-1",
				Name:                 "company-1",
				ProfilePictureBase64: "profile-picture-1-content",
			},
			upsertResponse: &models.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
				ProfilePicture:   "profile-picture-1",
			},
			expect: &linkedin_data_pb.Company{
				PublicIdentifier:  "public-identifier-1",
				Name:              "company-1",
				ProfilePictureUrl: "profile-picture-1",
			},
		},
		{
			name: "InvalidUpsertCompany",
			in: &linkedin_data_pb.UpsertCompanyRequest{
				PublicIdentifier: "public-identifier-2",
				Name:             "company-2",
			},
			upsertErr:  services.ErrInvalidUpsertCompany,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.UpsertCompanyRequest{
				PublicIdentifier: "public-identifier-3",
				Name:             "company-3",
			},
			upsertErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpsertCompanyService(t)
			service.
				On("Exec", context.TODO(), tt.in.PublicIdentifier, &models.UpsertCompany{
					Name:           tt.in.Name,
					ProfilePicture: tt.in.ProfilePictureBase64,
				}).
				Return(tt.upsertResponse, tt.upsertErr)

			handler := handlers.NewUpsertCompany(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.UpsertCompany(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
