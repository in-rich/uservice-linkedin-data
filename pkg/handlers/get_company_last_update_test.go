package handlers_test

import (
	"context"
	"errors"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	servicesmocks "github.com/in-rich/uservice-linkedin-data/pkg/services/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGetCompanyLastUpdateData(t *testing.T) {
	testData := []struct {
		name string

		in *linkedin_data_pb.GetCompanyLastUpdateRequest

		getResponse *time.Time
		getErr      error

		expect     *linkedin_data_pb.GetCompanyLastUpdateResponse
		expectCode codes.Code
	}{
		{
			name: "GetCompanyLastUpdateHandler",
			in: &linkedin_data_pb.GetCompanyLastUpdateRequest{
				PublicIdentifier: "public-identifier-1",
			},
			getResponse: lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			expect: &linkedin_data_pb.GetCompanyLastUpdateResponse{
				UpdatedAt: timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectCode: codes.OK,
		},
		{
			name: "CompanyLastUpdateNotFound",
			in: &linkedin_data_pb.GetCompanyLastUpdateRequest{
				PublicIdentifier: "public-identifier-2",
			},
			getErr:     dao.ErrCompanyNotFound,
			expectCode: codes.NotFound,
		},
		{
			name: "InternalError",
			in: &linkedin_data_pb.GetCompanyLastUpdateRequest{
				PublicIdentifier: "public-identifier-3",
			},
			getErr:     errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetCompanyLastUpdateService(t)
			service.On("Exec", context.TODO(), tt.in.PublicIdentifier).Return(tt.getResponse, tt.getErr)

			handler := handlers.NewGetCompanyLastUpdate(service)
			resp, err := handler.GetCompanyLastUpdate(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
