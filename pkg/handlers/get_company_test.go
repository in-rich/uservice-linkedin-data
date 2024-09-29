package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	servicesmocks "github.com/in-rich/uservice-linkedin-data/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestGetCompanyData(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.GetCompanyRequest

		getResponse *models.Company
		getErr      error

		expect     *linkedin_data_pb.Company
		expectCode codes.Code
	}{
		{
			name: "GetCompanyHandler",
			in: &linkedin_data_pb.GetCompanyRequest{
				PublicIdentifier: "public-identifier-1",
			},
			getResponse: &models.Company{
				PublicIdentifier: "public-identifier-1",
				Name:             "first-name-1",
				ProfilePicture:   "profile-picture-1",
			},
			expect: &linkedin_data_pb.Company{
				PublicIdentifier:  "public-identifier-1",
				Name:              "first-name-1",
				ProfilePictureUrl: "profile-picture-1",
			},
			expectCode: codes.OK,
		},
		{
			name: "CompanyNotFound",
			in: &linkedin_data_pb.GetCompanyRequest{
				PublicIdentifier: "public-identifier-2",
			},
			getErr:     dao.ErrCompanyNotFound,
			expectCode: codes.NotFound,
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.GetCompanyRequest{
				PublicIdentifier: "public-identifier-3",
			},
			getErr:     errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetCompanyService(t)
			service.On("Exec", context.TODO(), tt.in.PublicIdentifier).Return(tt.getResponse, tt.getErr)

			handler := handlers.NewGetCompany(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.GetCompany(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
