package services_test

import (
	"context"
	daomocks "github.com/in-rich/uservice-linkedin-data/pkg/dao/mocks"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetCompanyLastUpdate(t *testing.T) {
	testData := []struct {
		name string

		publicIdentifier string

		getResponse *time.Time
		getErr      error

		expect    *time.Time
		expectErr error
	}{
		{
			name:                        "GetCompanyLastUpdate",
			publicIdentifier:            "public-identifier-1",
			getResponse:                 lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			expect:                      lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:             "GetError",
			publicIdentifier: "public-identifier-1",
			getErr:           FooErr,
			expectErr:        FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			getCompanyLastUpdate := daomocks.NewMockGetCompanyLastUpdateRepository(t)
			getCompanyLastUpdate.On("GetCompanyLastUpdate", context.TODO(), tt.publicIdentifier).Return(tt.getResponse, tt.getErr)

			service := services.NewGetCompanyLastUpdateService(getCompanyLastUpdate)
			res, err := service.Exec(context.TODO(), tt.publicIdentifier)

			require.Equal(t, tt.expect, res)
			require.Equal(t, tt.expectErr, err)

			getCompanyLastUpdate.AssertExpectations(t)
		})
	}
}
