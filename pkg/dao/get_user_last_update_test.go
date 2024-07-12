package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var getUserLastUpdateFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestGetUserLastUpdate(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		publicIdentifier string
		expect           *time.Time
		expectErr        error
	}{
		{
			name:             "GetUserLastUpdate",
			publicIdentifier: "public-identifier-1",
			expect:           lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:             "UserLastUpdateNotFound",
			publicIdentifier: "public-identifier-2",
			expectErr:        dao.ErrUserNotFound,
		},
	}

	stx := BeginTX(db, getUserLastUpdateFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewGetUserLastUpdateRepository(tx)
			userLastUpdate, err := repo.GetUserLastUpdate(context.TODO(), data.publicIdentifier)

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, userLastUpdate)
		})
	}
}
